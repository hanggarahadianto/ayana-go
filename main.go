package main

import (
	"ayana/db"
	"ayana/routes"
	utilsEnv "ayana/utils/env"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	configure, err := utilsEnv.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables ", err)
	}
	log.Println("Loaded CLIENT_ORIGIN:", configure.ClientOrigin)

	// Initialize database
	db.InitializeDb(&configure)

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	r := gin.Default()

	// Check if CLIENT_ORIGIN is empty
	clientOrigin := configure.ClientOrigin
	log.Println("Loaded CLIENT_ORIGIN from config:", clientOrigin)

	// Fallback: Try to load from os.Getenv if it's empty
	if clientOrigin == "" {
		clientOrigin = os.Getenv("CLIENT_ORIGIN")
		log.Println("Loaded CLIENT_ORIGIN from environment variable:", clientOrigin)
	}

	if clientOrigin == "" {
		log.Fatal("❌ CLIENT_ORIGIN is not set. Check your .env file or environment variables.")
	}

	// Initialize database
	db.InitializeDb(&configure)

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router

	// Apply CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://ayanagroup99.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	log.Fatal(r.Run("0.0.0.0:" + configure.ServerPort)) // ✅ Important for Docker

	log.Println("✅ CORS Middleware Applied Successfully!")

	// ************* Router
	routes.SetupAuthRouter(r)
	routes.SetupHomeRouter(r)
	routes.SetupReservationRouter(r)
	routes.SetupInfoRouter(r)
	routes.SetupProjectRouter(r)
	routes.SetupWeeklyProgressRouter(r)
	routes.SetupCashFlowRouter(r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to my Ayana application! 🚀",
		})
		log.Println("Welcome to my Ayana application! 🚀")
	})

	log.Println("🚀 Server running on port:", configure.ServerPort)
	log.Fatal(r.Run("0.0.0.0:" + configure.ServerPort)) // ✅ Works inside Docker

}
