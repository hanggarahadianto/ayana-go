package controllers

import (
	"ayana/service/hr"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPresence(c *gin.Context) {
	response, err := hr.GetPresenceService(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Gagal mengambil data presensi",
			"status":  "gagal",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Data presensi berhasil diambil",
		"status":  "sukses",
	})
}
