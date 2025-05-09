package controllers

import (
	"ayana/db"
	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetClusterByID(c *gin.Context) {
	id := c.Param("id") // ambil ID dari path parameter

	var cluster models.Cluster

	result := db.DB.First(&cluster, "id = ?", id)
	if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Cluster not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   cluster,
	})
}
