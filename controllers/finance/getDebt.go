package controller

import (
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOutstandingDebts(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}
	debtStatus := c.DefaultQuery("status", "")
	summaryOnlyStr := c.DefaultQuery("summary_only", "false")
	summaryOnly := summaryOnlyStr == "true"
	if summaryOnlyStr != "true" && summaryOnlyStr != "false" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter summary_only harus 'true' atau 'false'."})
		return
	}

	dateFilter, err := helper.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}

	pagination := helper.GetPagination(c)

	params := service.DebtFilterParams{
		CompanyID:  companyID.String(),
		Pagination: pagination,
		DateFilter: dateFilter,
		DebtStatus: debtStatus,

		SummaryOnly: summaryOnly,
	}

	data, totalDebt, total, err := service.GetDebtsFromJournalLines(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data aset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"debtList":   data,
			"total_debt": totalDebt,
			"page":       pagination.Page,
			"limit":      pagination.Limit,
			"total":      total,
		},
		"message": "Hutang berhasil diambil",
		"status":  "sukses",
	})
}
