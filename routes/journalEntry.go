package routes

import (
	journalEntryController "ayana/controllers/journalEntry"
	"ayana/utils/handler"

	"github.com/gin-gonic/gin"
)

func SetupJournalEntryRouter(r *gin.Engine) {
	journalEntry := r.Group("/journal-entry")
	{
		journalEntry.GET("/get", journalEntryController.GetJournalEntriesByCategory)
		journalEntry.POST("/post", journalEntryController.CreateJournalEntry)
		journalEntry.DELETE("/delete", journalEntryController.DeleteJournalEntries)
		journalEntry.PUT("/update", journalEntryController.UpdateJournalEntry)
		journalEntry.POST("/reversed-post", journalEntryController.CreateReversedJournalEntry)

		journalEntry.POST("/download-transaction", handler.DownloadJournalTransactionHandler)

	}
}
