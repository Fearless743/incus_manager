package main

import (
	"log"
	"os"

	"incus-manager/internal/config"
	"incus-manager/internal/handler"
	"incus-manager/internal/service"
)

func main() {
	cfg := config.Load()
	
	db, err := database.InitDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	incusService := service.NewIncusService(cfg.IncusURL, cfg.IncusCert)
	authService := service.NewAuthService(db)
	userService := service.NewUserService(db)
	hostService := service.NewHostService(db, incusService)
	instanceService := service.NewInstanceService(db, incusService)
	sharedService := service.NewSharedService(db, incusService)

	h := handler.NewHandler(authService, userService, hostService, instanceService, sharedService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := h.Start(": " + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
