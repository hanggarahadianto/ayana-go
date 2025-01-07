package db

import (
	"ayana/models"
	utilsEnv "ayana/utils/env"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitializeDb(config *utilsEnv.Config) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,                               // Disable prepared statements
		Logger:      logger.Default.LogMode(logger.Info), // Log SQL queries
	})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	// Set up connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to configure database connection pool")
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	fmt.Println("Connected to database")

	// Ensure required extensions exist
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Migrate models
	db.AutoMigrate(
		&models.Home{},
		&models.Image{},
		&models.Reservation{},
		&models.Marketing{},
		&models.Info{},
		&models.User{},
	)

	DB = db
}
