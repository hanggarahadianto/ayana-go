package controllers

import (
	"ayana/db"
	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeletePresenceRule handles DELETE /presence-rules/:id
func DeletePresenceRule(c *gin.Context) {
	id := c.Param("id")
	ruleID, err := uuid.Parse(id)
	if err != nil {
		log.Println("🔴 Invalid UUID:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid.",
		})
		return
	}

	var rule models.PresenceRule
	if err := db.DB.First(&rule, "id = ?", ruleID).Error; err != nil {
		log.Println("🔴 Rule tidak ditemukan:", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Aturan presensi tidak ditemukan.",
		})
		return
	}

	if err := db.DB.Delete(&rule).Error; err != nil {
		log.Println("🔴 Gagal menghapus:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus aturan presensi.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Aturan presensi berhasil dihapus.",
	})
}
