package main

import (
	"bookcabin-voucher/config"
	"bookcabin-voucher/infrastructure/persistent"
	"bookcabin-voucher/internal/api/handler"
	http "bookcabin-voucher/internal/api/v1"
	"bookcabin-voucher/internal/middleware"
	"bookcabin-voucher/internal/migration"
	"bookcabin-voucher/internal/service"
	"bookcabin-voucher/internal/usecase"
	"bookcabin-voucher/internal/validation"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func main() {
	// Load environment variables
	cfg := config.LoadConfig()

	// Connect to SQLite
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Run DB migration
	migration.Migrate(db)

	// Init dependencies
	repo := persistent.NewFlightRepository(db)
	seatGenerator := service.NewSeatAllocator(cfg.SeatLayoutPath)
	u := usecase.NewFlightUsecase(repo, seatGenerator)
	h := handler.NewFlightHandler(u)

	// Setup Gin
	r := gin.Default()
	r.Use(middleware.CORSMiddleware([]string{cfg.FrontendURL}))
	r.Use(middleware.RecoveryMiddleware())

	//add custom validation

	validation.RegisterValidators()

	// Register routes
	http.RegisterRoutes(r, h)

	// Run server
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
