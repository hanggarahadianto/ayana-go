package controllers

import (
	"ayana/service"
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

	pagination := helper.GetPagination(c)

	search := c.Query("search") // optional kalau kamu pakai

	params := service.ProjectFilterParams{
		CompanyID:  companyID.String(),
		Pagination: pagination,
		Search:     search,
	}

	projects, totalProject, total, err := service.GetProjects(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"projectList":   projects,
			"total_project": totalProject,
			"page":          pagination.Page,
			"limit":         pagination.Limit,
			"total":         total,
		},
		"message": "Data project berhasil diambil",
		"status":  "sukses",
	})
}
