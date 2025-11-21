package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	url := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	return db, nil
}
