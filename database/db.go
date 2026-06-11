package database

import (
	"ai-incident-manager/models"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "incidents.db"
	}
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
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
