package routes

import (
	financeController "ayana/controllers/finance"

	"github.com/gin-gonic/gin"
)

func SetupFinanceRouter(r *gin.Engine) {
	finance := r.Group("/finance")
	{
		finance.GET("/get-outstanding-debt", financeController.GetOutstandingDebts)
		finance.GET("/get-expense-summary", financeController.GetExpensesSummary)
		finance.GET("/get-asset-summary", financeController.GetAssetSummary)
		finance.GET("/get-equity-summary", financeController.GetEquitySummary)
		finance.GET("/get-revenue-summary", financeController.GetRevenueSummary)

		// journalEntry.POST("/post", journalEntryController.CreateJournalEntry)
	}
}
