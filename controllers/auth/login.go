package controllers

import (
	"ayana/db"
	"ayana/models"
	"fmt"
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

	// Bind JSON input
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	// Cetak username dan password ke console

	// Cari user berdasarkan username
	var user models.User
	result := db.DB.First(&user, "username = ?", loginData.Username)
	if result.Error != nil {
		fmt.Println("User not found")
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Username not found",
		})
		return
	}

	// Validasi password (sebaiknya gunakan bcrypt.CompareHashAndPassword di produksi)
	if err := utilsAuth.VerifiedPassword(user.Password, loginData.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  false,
			"message": "Wrong password",
		})
		return
	}
	// Buat token JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	config, _ := utilsEnv.LoadConfig(".")

	fmt.Println("JWT ENV", []byte(config.JWTSecret))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to generate token",
		})
		return
	}

	// Berikan token sebagai respons
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Login success",
		"data": gin.H{
			"user":  user,
			"token": tokenString,
		},
	})
}
