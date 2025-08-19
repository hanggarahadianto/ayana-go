package controllers

import (
	lib "ayana/lib"
	"ayana/models"
	project "ayana/service/project"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProject(c *gin.Context) {
	companyIDStr := c.DefaultQuery("company_id", "")

	// Validasi Company ID
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}
	dateFilter, err := lib.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}

	pagination := lib.GetPagination(c)

	search := c.Query("search") // optional kalau kamu pakai

	params := project.ProjectFilterParams{
		CompanyID:  companyID.String(),
		Pagination: pagination,
		DateFilter: dateFilter,
		Search:     search,
	}

	projects, totalProject, total, err := project.GetProjects(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var enriched []models.ProjectWithStatus
	for _, proj := range projects {
		enriched = append(enriched, project.EnrichProjectStatus(proj))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"projectList":   enriched,
			"total_project": totalProject,
			"page":          pagination.Page,
			"limit":         pagination.Limit,
			"total":         total,
		},
		"message": "Data project berhasil diambil",
		"status":  "sukses",
	})

}
