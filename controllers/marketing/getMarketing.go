package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMarketing(c *gin.Context) {
	var marketingList []models.Marketing

	result := db.DB.Debug().Order("created_at desc, updated_at desc").Find(&marketingList)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   marketingList,
	})

}
