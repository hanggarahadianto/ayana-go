package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UpdateCluster(c *gin.Context) {
	var input models.Cluster
	id := c.Param("id")

	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cluster models.Cluster

	// Cari cluster berdasarkan ID
	if err := db.DB.First(&cluster, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster tidak ditemukan"})
		return
	}

	// Update field yang diizinkan
	cluster.Name = input.Name
	cluster.Location = input.Location
	cluster.Square = input.Square
	cluster.Price = input.Price
	cluster.Quantity = input.Quantity
	cluster.Status = input.Status
	cluster.Sequence = input.Sequence
	cluster.Maps = input.Maps
	cluster.UpdatedAt = time.Now()

	for i := range input.NearBies {
		input.NearBies[i].ID = uuid.New()
		input.NearBies[i].ClusterID = cluster.ID
	}
	cluster.NearBies = input.NearBies

	// Simpan perubahan
	if err := db.DB.Save(&cluster).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate cluster"})
		return
	}

	c.JSON(http.StatusOK, cluster)
}
