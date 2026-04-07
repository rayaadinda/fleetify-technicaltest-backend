package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/rayaadinda/fleetify-technicaltest-backend/internal/models"
)

type ItemHandler struct {
	db *gorm.DB
}

func NewItemHandler(db *gorm.DB) *ItemHandler {
	return &ItemHandler{db: db}
}

func (h *ItemHandler) GetItems(c *fiber.Ctx) error {
	codeQuery := strings.TrimSpace(c.Query("code"))

	query := h.db.Model(&models.Item{}).Order("code ASC").Limit(10)
	if codeQuery != "" {
		query = query.Where("code ILIKE ?", "%"+codeQuery+"%")
	}

	var items []models.Item
	if err := query.Find(&items).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to fetch items"})
	}

	return c.JSON(fiber.Map{"data": items})
}
