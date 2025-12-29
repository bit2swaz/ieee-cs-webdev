package database

import (
	"fmt"
	"log"
	"os"

	"github.com/bit2swaz/ieee-cs-webdev/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("Connected to Database successfully")

	// Enable UUID extension for Postgres
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Run Migrations
	log.Println("Running Migrations...")
	err = db.AutoMigrate(
		&models.Organization{},
		&models.User{},
		&models.Event{},
		&models.SubEvent{},
		&models.Ticket{},
	)
	if err != nil {
		log.Fatal("Migration Failed:  \n", err)
	}
	log.Println("Migrations Completed!")

	DB = db
}
