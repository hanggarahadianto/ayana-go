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
	// Connection string for PostgreSQL without invalid parameter
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
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

	// Ensure UUID extension for PostgreSQL (if needed)
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatalf("Failed to create UUID extension: %v", err)
	}

	fmt.Println("Connected to the database")

	// Auto-migrate models
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Company{},
		&models.Home{},
		&models.Info{},
		&models.NearBy{},
		&models.Reservation{},
		&models.Project{},
		&models.WeeklyProgress{},
		&models.Worker{},
		&models.Material{},
		&models.CashFlow{},
		&models.Goods{},
		// &models.Payout{},
		&models.Account{},
		&models.JournalEntry{},
		&models.JournalLine{},
		&models.TransactionCategory{},
		&models.Installment{},
	}

	for _, model := range modelsToMigrate {
		modelType := fmt.Sprintf("%T", model)

		// Check if the table exists for the current model
		if db.Migrator().HasTable(model) {
			fmt.Printf("Table for %s already exists, skipping auto-migration\n", modelType)
		} else {
			// If the table doesn't exist, create it
			fmt.Printf("Creating table for %s\n", modelType)
			err := db.Migrator().CreateTable(model)
			if err != nil {
				log.Printf("Failed to create table for %s: %v", modelType, err)
			}
		}
	}

	// Assign the DB instance to the global variable
	DB = db
}
