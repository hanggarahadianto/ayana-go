package controller

import (
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAssetSummary(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}

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

	assetType := c.DefaultQuery("asset_type", "")
	transactionType := c.DefaultQuery("transaction_type", "")
	pagination := helper.GetPagination(c)

	params := service.AssetFilterParams{
		CompanyID:       companyID.String(),
		Pagination:      pagination,
		DateFilter:      dateFilter,
		AssetType:       assetType,
		TransactionType: transactionType,
		SummaryOnly:     summaryOnly,
	}

	data, totalAsset, total, err := service.GetAssetsFromJournalLines(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data aset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"assetList":   data,
			"total_asset": totalAsset,
			"page":        pagination.Page,
			"limit":       pagination.Limit,
			"total":       total,
		},
		"message": "Aset berhasil diambil",
		"status":  "sukses",
	})
}
