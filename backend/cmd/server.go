package main

import (
	"log"
	"net/http"
	"os"
	"strings"

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

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"incus-manager"}`))
	})

	// WebSocket
	router.Handle("/ws", hub)

	// API routes
	auth := middleware.Authenticate(authService)
	router.HandleFunc("POST /api/login", h.Login)
	router.HandleFunc("POST /api/users", h.CreateUser)
	router.HandleFunc("POST /api/hosts", auth(h.AddHost))
	router.HandleFunc("GET /api/hosts", auth(h.GetHosts))
	router.HandleFunc("GET /api/instances", auth(h.GetInstances))
	router.HandleFunc("POST /api/instances", auth(h.CreateInstance))
	router.HandleFunc("DELETE /api/instances/", auth(h.DeleteInstance))
	router.HandleFunc("POST /api/instances/start/", auth(h.StartInstance))
	router.HandleFunc("POST /api/instances/stop/", auth(h.StopInstance))
	router.HandleFunc("POST /api/share", auth(h.ShareInstance))
	router.HandleFunc("DELETE /api/share/", auth(h.RevokeShare))
	router.HandleFunc("GET /api/instances/images", auth(h.GetImages))
	router.HandleFunc("GET /api/stats", auth(h.GetStats))

	// Static files
	router.HandleFunc("/", staticFileHandler())

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func staticFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && !strings.HasPrefix(r.URL.Path, "/assets/") &&
			!strings.HasSuffix(r.URL.Path, ".html") && !strings.HasSuffix(r.URL.Path, ".css") &&
			!strings.HasSuffix(r.URL.Path, ".js") && !strings.HasSuffix(r.URL.Path, ".svg") {
			http.NotFound(w, r)
			return
		}

		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		filePath := "/root/dist" + path
		data, err := os.ReadFile(filePath)
		if err != nil {
			indexData, err2 := os.ReadFile("/root/dist/index.html")
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
	lower := strings.ToLower(path)
	switch {
	case strings.HasSuffix(lower, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(lower, ".css"):
		return "text/css"
	case strings.HasSuffix(lower, ".js"):
		return "application/javascript"
	case strings.HasSuffix(lower, ".json"):
		return "application/json"
	case strings.HasSuffix(lower, ".png"):
		return "image/png"
	case strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(lower, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(lower, ".ico"):
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
