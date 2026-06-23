package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	incusService := service.NewIncusService(getIncusURL(), "", "")
	authService := service.NewAuthService(db, cfg.JWTSecret)
	userService := service.NewUserService(db)
	hostService := service.NewHostService(db, incusService)
	ipManager := service.NewIPManager(db)
	instanceService := service.NewInstanceService(db, incusService, ipManager)
	sharedService := service.NewSharedService(db, incusService)
	hub := websocket.NewHub()

	go hub.Run()

	h := handler.NewHandler(authService, userService, hostService, instanceService, sharedService, ipManager)

	router := http.NewServeMux()

	// API routes
	apiHandler := middleware.CORSMiddleware()(middleware.LoggingMiddleware(h.RegisterRoutes()))
	router.Handle("/api/", apiHandler)

	// WebSocket
	router.Handle("/ws", hub)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"incus-manager"}`))
	})

	// Static files - serve frontend
	router.Handle("GET /", staticFileHandler())

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

		cleanPath := strings.TrimPrefix(path, "/")
		filePath := filepath.Join("dist", cleanPath)

		data, err := os.ReadFile(filePath)
		if err != nil {
			// SPA fallback - serve index.html
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
	ext := strings.ToLower(filepath.Ext(path))
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

func getIncusURL() string {
	url := os.Getenv("INCUS_URL")
	if url == "" {
		if _, err := os.Stat("/var/run/incus/unix.sock"); err == nil {
			return "http://unix.socket"
		}
		url = "https://localhost:8443"
	}
	return url
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
