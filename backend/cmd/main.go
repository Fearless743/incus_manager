package main

import (
	"log"
	"net/http"
	"os"

	"incus-manager/internal/config"
	"incus-manager/internal/handler"
	"incus-manager/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := loadConfig()
	
	db, err := initDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	incusService := service.NewIncusService(cfg.IncusURL)
	authService := service.NewAuthService(db, cfg.JWTSecret)
	userService := service.NewUserService(db)
	hostService := service.NewHostService(db, incusService)
	instanceService := service.NewInstanceService(db, incusService)
	sharedService := service.NewSharedService(db, incusService)

	h := handler.NewHandler(authService, userService, hostService, instanceService, sharedService)

	router := h.RegisterRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func loadConfig() *config.Config {
	return &config.Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@localhost/incus_manager?sslmode=disable"),
		IncusURL:    getEnv("INCUS_URL", "https://localhost:8443"),
		JWTSecret:   getEnv("JWT_SECRET", "change-this-secret-in-production"),
	}
}

func initDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate models
	err = db.AutoMigrate()
	if err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
