package controller

import (
	lib "ayana/lib"
	service "ayana/service/marketing"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMarketerPerformanceHandler(c *gin.Context) {
	companyID := c.Query("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "company_id wajib diisi"})
		return
	}

	isAgentStr := c.Query("is_agent")
	var isAgent *bool
	if isAgentStr != "" {
		val := isAgentStr == "true"
		isAgent = &val
	}

	dateFilter, err := lib.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Format tanggal tidak valid"})
		return
	}

	result, err := service.GetMarketerPerformance(service.MarketerFilterParams{
		CompanyID:  companyID,
		IsAgent:    isAgent,
		DateFilter: dateFilter,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Gagal mengambil data performa marketer", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data performa marketer berhasil diambil",
		"data":    result,
	})
}
