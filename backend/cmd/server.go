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

	incusService := service.NewIncusService(cfg.IncusURL, "", "")
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
	router.Handle("/", middleware.CORSMiddleware()(middleware.LoggingMiddleware(h.RegisterRoutes())))
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"incus-manager"}`))
	})
	router.HandleFunc("/ws", hub.ServeHTTP)

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("WebSocket endpoint: ws://localhost:%s/ws", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
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
