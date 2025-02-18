package controllers

import (
	"ayana/db"
	middlewares "ayana/middlewares/token"
	"ayana/models"
	utilsAuth "ayana/utils/auth"
	utilsEnv "ayana/utils/env"

	"net/http"

	"github.com/gin-gonic/gin"
)

// Login handler for user authentication
func Login(c *gin.Context) {
	var loginData models.LoginData

	// Bind the JSON request data to loginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	// Find user by username
	var user models.User
	result := db.DB.First(&user, "username = ?", loginData.Username)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Username not found",
		})
		return
	}

	// Verify password
	if err := utilsAuth.VerifiedPassword(user.Password, loginData.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Wrong password",
		})
		return
	}

	// Generate JWT token
	config, _ := utilsEnv.LoadConfig(".")
	token, err := middlewares.GenerateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Error generating token",
			"error":   err.Error(),
		})
		return
	}

	// Set token as a secure HTTP-only cookie
	c.SetCookie("token", token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)

	// Respond with the user data and token
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Login success",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})
}
