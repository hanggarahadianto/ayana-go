package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeletePresenceRequest struct {
	IDs []string `json:"ids"`
}

func DeletePresence(c *gin.Context) {
	var req DeletePresenceRequest

	// 🔎 Validasi payload JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Format request tidak valid",
			"detail": err.Error(),
		})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Array ID tidak boleh kosong"})
		return
	}

	// 🧹 Proses delete dari database
	if err := db.DB.Where("id IN ?", req.IDs).Delete(&models.Presence{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Gagal menghapus data presensi",
			"detail": err.Error(),
		})
		return
	}

	// ✅ Berhasil
	c.JSON(http.StatusOK, gin.H{
		"status":  "sukses",
		"message": "Presensi berhasil dihapus",
		"deleted": len(req.IDs),
	})
}
