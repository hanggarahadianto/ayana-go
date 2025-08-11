package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	username, _ := c.Get("username")
	if username != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Hanya superadmin yang dapat mengakses data ini",
			"status":  "error",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var userList []models.User
	var total int64

	if err := db.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghitung total data",
		})
		return
	}

	if err := db.DB.Model(&models.User{}).
		Offset(offset).
		Limit(limit).
		Find(&userList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"userList":   userList,
			"status":     "success",
			"page":       page,
			"limit":      limit,
			"total_data": total,
			"total_page": (total + int64(limit) - 1) / int64(limit),
		},
		"message": "Berhasil mengambil data user",
		"status":  "success",
	})

}
