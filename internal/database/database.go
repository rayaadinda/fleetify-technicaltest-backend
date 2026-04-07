package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rayaadinda/fleetify-technicaltest-backend/internal/models"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Item{},
		&models.Invoice{},
		&models.InvoiceDetail{},
	)
}

func SeedItems(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Item{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	seedItems := []models.Item{
		{Code: "BRG-001", Name: "Ban Truk", Price: 1250000},
		{Code: "BRG-002", Name: "Oli Mesin", Price: 350000},
		{Code: "BRG-003", Name: "Filter Udara", Price: 150000},
		{Code: "BRG-004", Name: "Kampas Rem", Price: 475000},
		{Code: "BRG-005", Name: "Aki", Price: 980000},
	}

	return db.Create(&seedItems).Error
}
