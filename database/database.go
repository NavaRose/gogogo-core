package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func InitDatabaseWithoutEngine() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=%s password=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("TIME_ZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening database connection")
	}

	return db
}

func CloseDatabase(db *gorm.DB) {
	defer func() {
		sqlDB, _ := db.DB()
		if err := sqlDB.Close(); err != nil {
			log.Fatal("Failed to close database connection: ", err)
		}
	}()
}
