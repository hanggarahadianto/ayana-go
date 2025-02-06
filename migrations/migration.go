package migrations

import (
	"ayana/models"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func AddUnitToMaterial() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "2025020301", // Unique ID for this migration
		Migrate: func(tx *gorm.DB) error {
			// Check if the 'Unit' column exists in 'Material' table, if not, add it
			if !tx.Migrator().HasColumn(&models.Material{}, "Unit") {
				log.Println("Adding 'Unit' column to 'Material' table")
				if err := tx.Migrator().AddColumn(&models.Material{}, "Unit"); err != nil {
					log.Fatalf("Failed to add column 'Unit' to Material: %v", err)
					return err
				}
				log.Println("Successfully added column 'Unit' to Material table.")
			} else {
				log.Println("Column 'Unit' already exists in Material table.")
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// Drop the 'Unit' column if we need to rollback
			if tx.Migrator().HasColumn(&models.Material{}, "Unit") {
				log.Println("Rolling back migration and dropping 'Unit' column")
				if err := tx.Migrator().DropColumn(&models.Material{}, "Unit"); err != nil {
					log.Fatalf("Failed to drop column 'Unit' from Material: %v", err)
					return err
				}
				log.Println("Successfully dropped column 'Unit' from Material table.")
			}
			return nil
		},
	}
}
