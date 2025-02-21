package main

import (
	"ayana/db"
	"ayana/routes"
	utilsEnv "ayana/utils/env"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	configure, err := utilsEnv.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables ", err)
	}
	fmt.Println("Loaded CLIENT_ORIGIN:", configure.ClientOrigin)

	// Initialize database
	db.InitializeDb(&configure)

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	r := gin.Default()

	// Check if CLIENT_ORIGIN is empty
	if configure.ClientOrigin == "" {
		log.Fatal("❌ CLIENT_ORIGIN is not set. Check your .env file.")
	}

	// Apply CORS middleware before defining any routes
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{configure.ClientOrigin}, // ✅ Use the environment variable
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	fmt.Println("✅ CORS Middleware Applied Successfully!")

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
		fmt.Println("Welcome to my Ayana application! 🚀")
	})

	fmt.Println("🚀 Server running on port:", configure.ServerPort)
	log.Fatal(r.Run(":" + configure.ServerPort))
}
