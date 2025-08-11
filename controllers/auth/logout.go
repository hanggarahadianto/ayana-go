// Logout hanya memberi sinyal ke frontend untuk hapus token
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	// Set cookie kosong dengan expiry negatif → hapus cookie
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Logout success",
	})
}
