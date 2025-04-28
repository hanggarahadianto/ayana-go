package controller

import (
	"ayana/service"
	"ayana/utils/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAssetSummary mengembalikan ringkasan aset atau daftar aset
func GetAssetSummary(c *gin.Context) {
	// Validasi parameter wajib
	companyIDStr := c.Query("company_id")
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}

	summaryOnlyStr := c.DefaultQuery("summary_only", "false")

	summaryOnly := false
	if summaryOnlyStr == "true" {
		summaryOnly = true
	} else if summaryOnlyStr != "false" {
		// Menangani kasus di mana nilai selain "true" atau "false" diberikan
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter summary_only harus bernilai 'true' atau 'false'."})
		return
	}

	// Ambil filter tanggal
	dateFilter, err := helper.GetDateFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD."})
		return
	}

	// Ambil filter assetType dan transactionType
	assetType := c.DefaultQuery("asset_type", "") // Default ke empty string jika tidak ada
	transactionType := c.DefaultQuery("transaction_type", "")

	pagination := helper.GetPagination(c)

	params := service.AssetFilterParams{
		CompanyID:       companyID.String(),
		Pagination:      pagination,
		DateFilter:      dateFilter,
		AssetType:       assetType,       // Pasukkan assetType ke dalam params
		TransactionType: transactionType, // Pasukkan transactionType ke dalam params
		SummaryOnly:     summaryOnly,
	}

	data, totalAsset, total, err := service.GetAssets(params)
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
