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
	db.InitializeDb(&configure)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// ************* Router

	routes.SetupHomeRouter(r)
	routes.SetupReservationRouter(r)
	routes.SetupMarketingRouter(r)
	routes.SetupAuthRouter(r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to my Ayana application! 🚀",
		})
		fmt.Println("Welcome to my Ayana application! 🚀")
	})
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Replace with your allowed origins
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
	}))

	fmt.Println("running on server : " + configure.ServerPort)
	log.Fatal(r.Run(":" + configure.ServerPort))
}
