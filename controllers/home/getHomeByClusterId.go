package controllers

import (
	"ayana/db"
	"ayana/dto"
	lib "ayana/lib"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeListByClusterId(c *gin.Context) {
	clusterId := c.Param("cluster_id")
	pagination := lib.GetPagination(c)

	var homes []models.Home
	var total int64

	// Hitung total data
	db.DB.Model(&models.Home{}).Where("cluster_id = ?", clusterId).Count(&total)

	// Ambil data termasuk relasi Cluster dan NearBies (yang berasal dari Cluster)
	db.DB.
		Where("cluster_id = ?", clusterId).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Order("sequence asc").
		Preload("Cluster.NearBies").
		Preload("Cluster").
		Find(&homes)

	var homeResponses []dto.HomeByClusterResponse

	for _, home := range homes {
		// Mapping NearBies dari model ke DTO
		var nearBiesDTO []dto.NearBy
		for _, nearby := range home.Cluster.NearBies {
			nearBiesDTO = append(nearBiesDTO, dto.NearBy{
				ID:       nearby.ID.String(),
				Name:     nearby.Name,
				Distance: nearby.Distance,
			})
		}

		// Mapping Home + Cluster + NearBies
		homeResponse := dto.HomeByClusterResponse{
			ID:         home.ID.String(),
			Type:       home.Type,
			Title:      home.Title,
			Status:     home.Status,
			Maps:       home.Cluster.Maps,
			Quantity:   home.Quantity,
			Sequence:   home.Sequence,
			Bathroom:   home.Bathroom,
			Bedroom:    home.Bedroom,
			Content:    home.Content,
			Price:      home.Price,
			Square:     home.Square,
			StartPrice: home.StartPrice,
			Cluster: dto.ClusterResponse{
				ID:   home.Cluster.ID.String(),
				Name: home.Cluster.Name,
				Maps: home.Cluster.Maps,
			},
			NearBies: nearBiesDTO,
		}

		homeResponses = append(homeResponses, homeResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   homeResponses,
		"total":  total,
		"limit":  pagination.Limit,
		"offset": pagination.Offset,
	})
}
