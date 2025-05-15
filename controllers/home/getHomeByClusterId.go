package controllers

import (
	"ayana/db"
	"ayana/dto"
	"ayana/models"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomeListByClusterId(c *gin.Context) {
	clusterId := c.Param("cluster_id")
	pagination := helper.GetPagination(c)

	var homes []models.Home
	var total int64

	// Hitung total data
	db.DB.Model(&models.Home{}).Where("cluster_id = ?", clusterId).Count(&total)

	// Ambil data dengan NearBies saja, tanpa preload Images
	db.DB.
		Where("cluster_id = ?", clusterId).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Order("sequence asc").
		Preload("NearBies").
		Preload("Cluster"). // <= tambahkan ini
		Find(&homes)

	var homeResponses []dto.HomeByClusterResponse
	for _, home := range homes {
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
			NearBies: func() []dto.NearBy {
				var nearBies []dto.NearBy
				for _, nb := range home.NearBies {
					nearBies = append(nearBies, dto.NearBy{
						ID:       nb.ID.String(),
						Name:     nb.Name,
						Distance: nb.Distance,
					})
				}
				return nearBies
			}(), // Map models.NearBy to controllers.NearBy
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
