package database

import (
	"ai-incident-manager/models"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("incidents.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := DB.AutoMigrate(&models.Incident{}, &models.Message{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database initialized")
}
