package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/rayaadinda/fleetify-technicaltest-backend/internal/config"
	"github.com/rayaadinda/fleetify-technicaltest-backend/internal/handlers"
	"github.com/rayaadinda/fleetify-technicaltest-backend/internal/middleware"
)

func Setup(app *fiber.App, db *gorm.DB, cfg config.Config) {
	authHandler := handlers.NewAuthHandler(cfg.JWTSecret)
	itemHandler := handlers.NewItemHandler(db)
	invoiceHandler := handlers.NewInvoiceHandler(db)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	api := app.Group("/api")
	api.Post("/login", authHandler.Login)
	api.Get("/items", itemHandler.GetItems)
	api.Post("/invoices", middleware.Protected(cfg.JWTSecret), invoiceHandler.CreateInvoice)
}
