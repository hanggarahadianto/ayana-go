package utils

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsUsernameExists(c *gin.Context, username string) bool {
	var existingUser models.User
	db.DB.Where("username = ?", username).First(&existingUser)

	if existingUser.Username == username {
		// Username already exists
		c.JSON(http.StatusConflict, gin.H{
			"status":  false,
			"message": "Username already exists",
		})
		return true
	}

	return false
}
func FindUsername(c *gin.Context, username string) bool {
	var existingUser models.User

	db.DB.Where("username = ?", username).First(&existingUser)

	// if existingUser.Username == {
	// 	// Username already exists
	// 	c.JSON(http.StatusConflict, gin.H{
	// 		"status":  false,
	// 		"message": "Username already exists",
	// 	})
	// 	return true
	// }

	return false
}
