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
	// Connection string for Supabase
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable statement_cache_mode=none",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	// Connect to the database
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disable prepared statements in the driver
	}), &gorm.Config{
		PrepareStmt: false, // Disable prepared statements in GORM
		Logger:      logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to configure database connection pool: %v", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	// Ensure UUID extension for Supabase
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatalf("Failed to create UUID extension: %v", err)
	}

	fmt.Println("Connected to the database")

	// Auto-migrate models
	modelsToMigrate := []interface{}{
		&models.Home{},        // Base model should be first
		&models.Info{},        // Depends on Home
		&models.NearBy{},      // Depends on Info
		&models.Reservation{}, // Depends on Home
		&models.Project{},
	}

	for _, model := range modelsToMigrate {
		err := db.AutoMigrate(model)
		if err != nil {
			log.Printf("Failed to auto-migrate model %T: %v", model, err)
			continue // Continue with other migrations instead of fatal
		}
		fmt.Printf("Auto-migrated model %T successfully\n", model)
	}

	// Assign the DB instance to the global variable
	DB = db
}
