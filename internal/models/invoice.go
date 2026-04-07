package models

import "time"

type Invoice struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	InvoiceNumber   string          `gorm:"uniqueIndex;size:100;not null" json:"invoice_number"`
	SenderName      string          `gorm:"size:255;not null" json:"sender_name"`
	SenderAddress   string          `gorm:"type:text;not null" json:"sender_address"`
	ReceiverName    string          `gorm:"size:255;not null" json:"receiver_name"`
	ReceiverAddress string          `gorm:"type:text" json:"receiver_address"`
	TotalAmount     int64           `gorm:"not null" json:"total_amount"`
	CreatedBy       uint            `gorm:"not null" json:"created_by"`
	Details         []InvoiceDetail `json:"details"`
	CreatedAt       time.Time       `json:"created_at"`
}
