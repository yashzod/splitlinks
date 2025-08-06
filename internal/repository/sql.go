package repository

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/yashzod/splitlinks/internal/model" // update this import path as needed
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("yourapp.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to SQLite database:", err)
	}

	// Auto-migrate your models
	db.AutoMigrate(&model.Experiment{})
	db.AutoMigrate(&model.Variant{})
	db.AutoMigrate(&model.User{})

	return db
}
