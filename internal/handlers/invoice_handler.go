package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/rayaadinda/fleetify-technicaltest-backend/internal/models"
)

type InvoiceHandler struct {
	db *gorm.DB
}

type invoiceLineRequest struct {
	ItemCode string `json:"item_code"`
	Quantity int64  `json:"quantity"`
	Price    *int64 `json:"price,omitempty"`
	Total    *int64 `json:"total,omitempty"`
}

type createInvoiceRequest struct {
	SenderName      string               `json:"sender_name"`
	SenderAddress   string               `json:"sender_address"`
	ReceiverName    string               `json:"receiver_name"`
	ReceiverAddress string               `json:"receiver_address"`
	Items           []invoiceLineRequest `json:"items"`
}

func NewInvoiceHandler(db *gorm.DB) *InvoiceHandler {
	return &InvoiceHandler{db: db}
}

func (h *InvoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	var req createInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	if req.SenderName == "" || req.SenderAddress == "" || req.ReceiverName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "sender and receiver data are required"})
	}

	if len(req.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "at least one item is required"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid user context"})
	}

	invoiceNumber := generateInvoiceNumber()
	var totalAmount int64
	details := make([]models.InvoiceDetail, 0, len(req.Items))

	err := h.db.Transaction(func(tx *gorm.DB) error {
		for _, line := range req.Items {
			if line.ItemCode == "" || line.Quantity <= 0 {
				return errors.New("item_code and positive quantity are required")
			}

			var item models.Item
			if err := tx.Where("code = ?", line.ItemCode).First(&item).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("item with code %s not found", line.ItemCode)
				}
				return err
			}

			subtotal := item.Price * line.Quantity
			totalAmount += subtotal

			details = append(details, models.InvoiceDetail{
				ItemID:   item.ID,
				Quantity: line.Quantity,
				Price:    item.Price,
				Subtotal: subtotal,
			})
		}

		invoice := models.Invoice{
			InvoiceNumber:   invoiceNumber,
			SenderName:      req.SenderName,
			SenderAddress:   req.SenderAddress,
			ReceiverName:    req.ReceiverName,
			ReceiverAddress: req.ReceiverAddress,
			TotalAmount:     totalAmount,
			CreatedBy:       userID,
		}
		if err := tx.Create(&invoice).Error; err != nil {
			return err
		}

		for i := range details {
			details[i].InvoiceID = invoice.ID
		}

		if err := tx.Create(&details).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"invoice_number": invoiceNumber,
		"total_amount":   totalAmount,
	})
}

func generateInvoiceNumber() string {
	now := time.Now()
	return fmt.Sprintf("INV-%s-%04d", now.Format("20060102150405"), now.UnixNano()%10000)
}
