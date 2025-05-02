package routes

import (
	financeController "ayana/controllers/finance"

	"github.com/gin-gonic/gin"
)

func SetupFianceRouter(r *gin.Engine) {
	finance := r.Group("/finance")
	{
		finance.GET("/get-outstanding-debt", financeController.GetOutstandingDebts)
		finance.GET("/get-expense-summary", financeController.GetExpensesSummary)
		finance.GET("/get-asset-summary", financeController.GetAssetSummary)
		finance.GET("/get-available-cash", financeController.GetAvailableCashHandler)

		// journalEntry.POST("/post", journalEntryController.CreateJournalEntry)
	}
}
