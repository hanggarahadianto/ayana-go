package db

import (
	"ayana/models"
	"ayana/utils"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDb(config *utils.Config) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	fmt.Println("connected to database")

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	db.AutoMigrate(
		&models.Home{},
		&models.Image{},
		&models.Reservation{},
	)

	DB = db
}
