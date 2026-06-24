package main

import (
	"log"
	"net/http"
	"os"

	"incus-manager/internal/config"
	"incus-manager/internal/handler"
	"incus-manager/internal/middleware"
	"incus-manager/internal/service"
	"incus-manager/internal/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := initDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	authService := service.NewAuthService(db, cfg.JWTSecret)
	userService := service.NewUserService(db)
	incusFactory := service.NewIncusServiceFactory()
	hostService := service.NewHostService(db, incusFactory)
	ipManager := service.NewIPManager(db)
	instanceService := service.NewInstanceService(db, incusFactory, ipManager)
	sharedService := service.NewSharedService(db, incusFactory)
	hub := websocket.NewHub()

	go hub.Run()

	h := handler.NewHandler(authService, userService, hostService, instanceService, sharedService, ipManager)

	router := http.NewServeMux()

	// Health check (public, no middleware)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"incus-manager"}`))
	})

	// Static files - serve frontend (public, no middleware)
	router.Handle("GET /", staticFileHandler())

	// API routes (include /api/ prefix internally)
	apiHandler := middleware.CORSMiddleware()(middleware.LoggingMiddleware(h.RegisterRoutes()))
	router.Handle("/", apiHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Frontend: http://localhost:%s", port)
	log.Printf("API: http://localhost:%s/api", port)
	log.Printf("Health: http://localhost:%s/health", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func staticFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		cleanPath := path[1:]
		filePath := "dist/" + cleanPath

		data, err := os.ReadFile(filePath)
		if err != nil {
			indexData, err2 := os.ReadFile("dist/index.html")
			if err2 != nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(indexData)
			return
		}

		w.Header().Set("Content-Type", getContentType(path))
		w.Write(data)
	}
}

func getContentType(path string) string {
	ext := path[len(path)-4:]
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

func initDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
