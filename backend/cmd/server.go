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
	"incus-manager/internal/model"
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

	// Auto migrate all models
	if err := db.AutoMigrate(&model.User{}, &model.Host{}, &model.Instance{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
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
	router.HandleFunc("/api/login", h.Login)
	router.HandleFunc("/api/users", h.CreateUser)
	router.HandleFunc("/api/hosts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			auth(h.AddHost)(w, r)
		default:
			auth(h.GetHosts)(w, r)
		}
	})
	router.HandleFunc("/api/hosts/test", auth(h.TestHost))
	router.HandleFunc("/api/hosts/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/hosts/")
		if strings.Contains(path, "/test") {
			return
		}
		switch r.Method {
		case http.MethodPut:
			auth(h.UpdateHost)(w, r)
		case http.MethodDelete:
			auth(h.DeleteHost)(w, r)
		default:
			auth(h.GetHosts)(w, r)
		}
	})
	router.HandleFunc("/api/instances", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			auth(h.CreateInstance)(w, r)
		default:
			auth(h.GetInstances)(w, r)
		}
	})
	router.HandleFunc("/api/instances/", auth(h.DeleteInstance))
	router.HandleFunc("/api/instances/start/", auth(h.StartInstance))
	router.HandleFunc("/api/instances/stop/", auth(h.StopInstance))
	router.HandleFunc("/api/share", auth(h.ShareInstance))
	router.HandleFunc("/api/share/", auth(h.RevokeShare))
	router.HandleFunc("/api/instances/images", auth(h.GetImages))
	router.HandleFunc("/api/stats", auth(h.GetStats))

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
		if r.URL.Path == "/" {
			r.URL.Path = "/index.html"
		}

		staticExtensions := map[string]bool{
			".html": true, ".css": true, ".js": true, ".json": true,
			".png": true, ".jpg": true, ".jpeg": true, ".svg": true, ".ico": true, ".webp": true, ".woff": true, ".woff2": true, ".ttf": true, ".eot": true,
		}
		_, ext := staticExtensions[strings.ToLower(filepath.Ext(r.URL.Path))]

		if ext || strings.HasPrefix(r.URL.Path, "/assets/") {
			filePath := "/root/dist" + r.URL.Path
			data, err := os.ReadFile(filePath)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", getContentType(r.URL.Path))
			w.Write(data)
			return
		}

		indexData, err := os.ReadFile("/root/dist/index.html")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexData)
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
