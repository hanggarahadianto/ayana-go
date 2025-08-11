package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	utilsAuth "ayana/utils/auth"
	utilsEnv "ayana/utils/env"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type LoginData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {
	var loginData LoginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	var user models.User
	if err := db.DB.First(&user, "username = ?", loginData.Username).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Username not found",
		})
		return
	}

	if err := utilsAuth.VerifiedPassword(user.Password, loginData.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Wrong password",
		})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	config, _ := utilsEnv.LoadConfig(".")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to generate token",
		})
		return
	}

	// Simpan token di cookie (HttpOnly)
	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Login success",
		"data": gin.H{
			"user": user,
		},
	})
}
