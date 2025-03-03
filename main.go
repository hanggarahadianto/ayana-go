package main

import (
	"ayana/db"
	"ayana/routes"
	utilsEnv "ayana/utils/env"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func main() {
	// Setup logging
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	// Load environment variables
	configure, err := utilsEnv.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables: ", err)
	}

	// Load CLIENT_ORIGIN from config or environment variable
	clientOrigin := configure.ClientOrigin
	if clientOrigin == "" {
		clientOrigin = os.Getenv("CLIENT_ORIGIN")
	}
	if clientOrigin == "" {
		log.Fatal("❌ CLIENT_ORIGIN is not set. Check your .env file or environment variables.")
	}

	log.Println("Loaded CLIENT_ORIGIN:", clientOrigin)

	// Initialize database
	db.InitializeDb(&configure)

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	r := gin.Default()

	// Apply CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{clientOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true, // Set to true for debugging CORS issues
	})

	log.Println("✅ CORS Middleware Applied Successfully!")

	handler := c.Handler(r)

	// Setup routes
	routes.SetupAuthRouter(r)
	routes.SetupHomeRouter(r)
	routes.SetupReservationRouter(r)
	routes.SetupInfoRouter(r)
	routes.SetupProjectRouter(r)
	routes.SetupWeeklyProgressRouter(r)
	routes.SetupCashFlowRouter(r)

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to my Ayana application! 🚀",
		})
		log.Println("Welcome to my Ayana application! 🚀")
	})

	// Start server
	serverAddr := "0.0.0.0:" + configure.ServerPort
	log.Println("🚀 Server running on port:", configure.ServerPort)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}
