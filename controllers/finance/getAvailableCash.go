package controller

import (
	"ayana/service"
	"ayana/utils/helper"

	"github.com/gin-gonic/gin"
)

func GetAvailableCashHandler(c *gin.Context) {
	companyIDStr := c.Query("company_id")
	companyID, valid := helper.ValidateAndParseCompanyID(companyIDStr, c)
	if !valid {
		return
	}

	totalCashIn, totalCashOut, availableCash, err := service.GetAvailableCash(companyID.String())
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get available cash"})
		return
	}

	c.JSON(200, gin.H{
		"status": "sukses",
		"data": gin.H{
			"total_cash_in":  totalCashIn,
			"total_cash_out": totalCashOut,
			"available_cash": availableCash,
		},
	})
}
