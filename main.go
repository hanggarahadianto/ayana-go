package main

import (
	"ayana/db"
	"ayana/routes"
	utilsEnv "ayana/utils/env"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	configure, err := utilsEnv.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables ", err)
	}
	db.InitializeDb(&configure)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Apply CORS middleware before defining any routes
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},                                     // Allow requests from http://localhost:3000
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                   // Include OPTIONS for preflight
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, // Add Authorization
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Cache preflight response for 12 hours
	}))

	// Handle OPTIONS requests globally
	r.OPTIONS("/*any", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(http.StatusNoContent) // 204 No Content
	})

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

	fmt.Println("running on server : " + configure.ServerPort)
	log.Fatal(r.Run(":" + configure.ServerPort))
}
