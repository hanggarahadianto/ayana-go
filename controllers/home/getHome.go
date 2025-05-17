package controllers

import (
	"ayana/db"
	"ayana/dto"
	"ayana/utils/helper"

	"ayana/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHomes(c *gin.Context) {
	pagination := helper.GetPagination(c)

	// Validasi parameter pagination
	if !helper.ValidatePagination(pagination, c) {
		return
	}

	// Ambil filter status dari query
	status := c.Query("status")

	var homeList []models.Home
	var total int64

	// Mulai query DB
	query := db.DB.Model(&models.Home{})

	// Filter berdasarkan status jika ada
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Hitung total data dengan filter
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Count error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghitung data",
		})
		return
	}

	// Ambil data dengan limit, offset, dan order
	result := query.
		Preload("Cluster").
		Preload("Cluster.NearBies").
		Order("sequence asc").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&homeList)

	if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}

	var response []dto.HomeByClusterResponse

	for _, h := range homeList {
		var nearBiesDTO []dto.NearBy
		var clusterResp dto.ClusterResponse
		var maps, location string

		if h.Cluster != nil {
			maps = h.Cluster.Maps
			location = h.Cluster.Location

			for _, n := range h.Cluster.NearBies {
				nearBiesDTO = append(nearBiesDTO, dto.NearBy{
					ID:       n.ID.String(),
					Name:     n.Name,
					Distance: n.Distance,
				})
			}

			clusterResp = dto.ClusterResponse{
				ID:   h.Cluster.ID.String(),
				Name: h.Cluster.Name,
				Maps: h.Cluster.Maps,
			}
		}

		resp := dto.HomeByClusterResponse{
			ID:         h.ID.String(),
			Title:      h.Title,
			Type:       h.Type,
			Content:    h.Content,
			Maps:       maps,
			Location:   location,
			Price:      h.Price,
			Status:     h.Status,
			Quantity:   h.Quantity,
			Sequence:   h.Sequence,
			Square:     h.Square,
			Bathroom:   h.Bathroom,
			Bedroom:    h.Bedroom,
			StartPrice: h.StartPrice,
			Cluster:    clusterResp,
			NearBies:   nearBiesDTO,
		}
		response = append(response, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       response,
		"page":       pagination.Page,
		"limit":      pagination.Limit,
		"total_data": total,
		"total_page": (total + int64(pagination.Limit) - 1) / int64(pagination.Limit),
		"status":     "success",
	})

}
