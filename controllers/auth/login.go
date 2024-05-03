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

func Login(c *gin.Context) {

	var loginData models.LoginData

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// =================== find username
	var user models.User
	result := db.DB.First(&user, "username = ?", (loginData.Username))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Username not found"})
		return
	}

	if err := utilsAuth.VerifiedPassword(user.Password, loginData.Password); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": "wrong password"})
		return
	}
	config, _ := utilsEnv.LoadConfig(".")
	token, err := middlewares.GenerateToken(
		config.AccessTokenExpiresIn,
		user.ID,
		config.AccessTokenPrivateKey,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error()})
		return
	}
	c.SetCookie("token", token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Login Success",
		"data":    gin.H{"payload": user, "token": token},
	})

}
