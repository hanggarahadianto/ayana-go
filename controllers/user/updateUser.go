package controllers

import (
	"ayana/db"
	"ayana/models"
	utilsAuth "ayana/utils/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateUser(c *gin.Context) {
	// Cek role yang login
	username, _ := c.Get("username")
	if username != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Hanya superadmin yang dapat mengupdate data user",
			"status":  "error",
		})
		return
	}

	// Ambil ID user dari parameter URL
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ID user wajib diisi",
			"status":  "error",
		})
		return
	}

	// Cari user yang mau diupdate
	var user models.User
	if err := db.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User tidak ditemukan",
			"status":  "error",
		})
		return
	}

	// Bind JSON dari request body
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Data request tidak valid",
			"status":  "error",
		})
		return
	}

	// Update field
	if input.Username != "" {
		user.Username = input.Username
	}
	if input.Password != "" {
		// Hash password jika perlu
		hashed, _ := utilsAuth.HashPassword(input.Password) // Pastikan kamu punya fungsi hashing
		user.Password = hashed
	}
	if input.Role != "" {
		user.Role = input.Role
	}

	// Simpan perubahan
	if err := db.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengupdate data user",
			"status":  "error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "Berhasil mengupdate data user",
		"status":  "success",
	})
}
