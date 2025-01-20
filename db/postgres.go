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
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable statement_cache_mode=none",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disable prepared statements in the driver
	}), &gorm.Config{
		PrepareStmt: false, // Disable prepared statements in GORM
		Logger:      logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to configure database connection pool")
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	fmt.Println("Connected to database")

	// db.Exec("SET statement_cache_mode = 'none';") // Disable caching at runtime
	// db.Exec("DISCARD ALL;")                       // Clear existing cache

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// modelsToCreate := []interface{}{
	// 	&models.User{},
	// 	&models.Home{},
	// 	&models.Image{},
	// 	&models.Reservation{},
	// 	&models.Marketing{},
	// 	&models.Info{},
	// 	&models.NearBy{},
	// }

	// for _, model := range modelsToCreate {
	// 	err := db.Migrator().CreateTable(model)
	// 	if err != nil {
	// 		log.Printf("Failed to create table for model %T: %v", model, err)
	// 	} else {
	// 		fmt.Printf("Created table for model %T\n", model)
	// 	}
	// }

	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Home{},
		&models.Image{},
		&models.Reservation{},
		&models.Info{},
		&models.NearBy{},
	}

	for _, model := range modelsToMigrate {
		err := db.AutoMigrate(model)
		if err != nil {
			log.Printf("Failed to auto-migrate model %T: %v", model, err)
		} else {
			fmt.Printf("Auto-migrated model %T\n", model)
		}
	}

	fmt.Println("Connected to the database and migrated successfully.")
	DB = db

}
