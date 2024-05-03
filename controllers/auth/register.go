package controllers

import (
	"ayana/db"
	"ayana/models"
	utilsAuth "ayana/utils/auth"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {

	var registerData models.RegisterData

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if utilsAuth.IsUsernameExists(c, registerData.Username) {
		return
	}

	hashPassword, err := utilsAuth.HashPassword(registerData.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "password failed to hashed",
			"message": err.Error(),
		})

	}

	if registerData.Password != registerData.PasswordConfirm {
		errorResponse := gin.H{
			"status":  "failed",
			"message": "Passoword do not match",
		}
		c.JSON(http.StatusBadRequest, errorResponse)

		errorMessage, err := json.MarshalIndent(errorResponse, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
		}
		fmt.Printf("response:\n%s\n", string(errorMessage))

		return
	}

	now := time.Now()
	newUser := models.User{
		Username: registerData.Username,
		Password: hashPassword,
		Role:     registerData.Role,

		CreatedAt: now,
		UpdatedAt: now,
	}

	db.DB.Debug().Create(&newUser)

	c.JSON(http.StatusCreated, gin.H{
		"status": true,
		"data":   newUser,
	})

	fmt.Println(newUser)

}
