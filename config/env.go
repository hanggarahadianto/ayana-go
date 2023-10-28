package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file")
	}
}
