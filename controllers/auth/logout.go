// Logout hanya memberi sinyal ke frontend untuk hapus token
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	// Bisa tambahkan logging user logout di sini kalau mau
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Logout success. Token should be removed from the client.",
	})
}
