package controllers

import (
	"ayana/db"
	"ayana/utils/helper"

	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeByClusterResponse struct {
	ID       string          `json:"id"`
	Title    string          `json:"title"`
	Content  string          `json:"content"`
	Price    float64         `json:"price"`
	Square   float64         `json:"square"`
	Cluster  ClusterResponse `json:"cluster"`
	NearBies []NearBy        `json:"near_bies"`
}

type ClusterResponse struct {
	Location string `json:"location"`
	Maps     string `json:"maps"`
}

type NearBy struct {
	Name     string `json:"name"`
	Distance string `json:"distance"`
}

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

	var homeResponses []HomeByClusterResponse
	for _, home := range homes {
		homeResponse := HomeByClusterResponse{
			ID:      home.ID.String(),
			Title:   home.Title,
			Content: home.Content,
			Price:   home.Price,
			Square:  home.Square,
			Cluster: ClusterResponse{
				Location: home.Cluster.Location,
				Maps:     home.Cluster.Maps,
			},
			NearBies: func() []NearBy {
				var nearBies []NearBy
				for _, nb := range home.NearBies {
					nearBies = append(nearBies, NearBy{
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
