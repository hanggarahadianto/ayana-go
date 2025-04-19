package routes

import (
	financeController "ayana/controllers/finance"

	"github.com/gin-gonic/gin"
)

func SetupFianceRouter(r *gin.Engine) {
	finance := r.Group("/finance")
	{
		finance.GET("/get-outstanding-debt", financeController.GetOutstandingDebts)

		// journalEntry.POST("/post", journalEntryController.CreateJournalEntry)
	}
}
