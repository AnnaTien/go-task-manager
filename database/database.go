package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes and returns a GORM database connection instance.
func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db
}
