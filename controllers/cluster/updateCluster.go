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
	id := c.Param("id")

	var input models.Cluster
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil cluster yang ingin diupdate
	var cluster models.Cluster
	if err := db.DB.Preload("NearBies").First(&cluster, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cluster tidak ditemukan"})
		return
	}

	// Update field utama
	cluster.Name = input.Name
	cluster.Location = input.Location
	cluster.Square = input.Square
	cluster.Price = input.Price
	cluster.Quantity = input.Quantity
	cluster.Status = input.Status
	cluster.Sequence = input.Sequence
	cluster.Maps = input.Maps
	cluster.UpdatedAt = time.Now()

	if err := db.DB.Save(&cluster).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate cluster"})
		return
	}

	// ================================
	// Proses NearBies: Create/Update/Delete
	// ================================
	// 1. Buat map ID NearBies dari input
	inputNearByMap := make(map[uuid.UUID]models.NearBy)
	for _, nb := range input.NearBies {
		if nb.ID != uuid.Nil {
			inputNearByMap[nb.ID] = nb
		}
	}

	// 2. Update atau Delete NearBies yang lama
	for _, existing := range cluster.NearBies {
		if inputNb, found := inputNearByMap[existing.ID]; found {
			// Update NearBy
			existing.Name = inputNb.Name
			existing.Distance = inputNb.Distance
			if err := db.DB.Save(&existing).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update NearBy"})
				return
			}
			delete(inputNearByMap, existing.ID) // tandai sudah diproses
		} else {
			// Hapus NearBy yang tidak ada lagi di input
			if err := db.DB.Delete(&existing).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal hapus NearBy lama"})
				return
			}
		}
	}

	// 3. Tambah NearBies baru (yang belum punya ID)
	for _, nb := range input.NearBies {
		if nb.ID == uuid.Nil {
			nb.ID = uuid.New()
			nb.ClusterID = cluster.ID
			if err := db.DB.Create(&nb).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan NearBy baru"})
				return
			}
		}
	}

	// Ambil data terbaru
	var updated models.Cluster
	if err := db.DB.Preload("NearBies").First(&updated, "id = ?", cluster.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil cluster terbaru"})
		return
	}

	c.JSON(http.StatusOK, updated)
}
