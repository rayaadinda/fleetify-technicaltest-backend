package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/rayaadinda/fleetify-technicaltest-backend/internal/config"
	"github.com/rayaadinda/fleetify-technicaltest-backend/internal/database"
	"github.com/rayaadinda/fleetify-technicaltest-backend/internal/routes"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	if err := database.SeedItems(db); err != nil {
		log.Fatalf("database seeding failed: %v", err)
	}

	app := fiber.New()
	app.Use(cors.New())

	routes.Setup(app, db, cfg)

	log.Printf("server listening on :%s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
