package controllers

import (
	"ayana/db"
	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCluster(c *gin.Context) {
	var clusterList []models.Cluster

	result := db.DB.
		Preload("NearBies").
		Order("sequence asc").Find(&clusterList)
	if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   clusterList,
	})

}
