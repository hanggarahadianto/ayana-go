package routes

import (
	journalEntryController "ayana/controllers/journalEntry"

	"github.com/gin-gonic/gin"
)

func SetupJournalEntryRouter(r *gin.Engine) {
	journalEntry := r.Group("/journal-entry")
	{
		journalEntry.GET("/get", journalEntryController.GetJournalEntriesByCategory)
		journalEntry.POST("/post", journalEntryController.CreateJournalEntry)
		journalEntry.DELETE("/delete/:id", journalEntryController.DeleteJournalEntry)
		journalEntry.PUT("/update", journalEntryController.UpdateJournalEntry)
		journalEntry.POST("/reversed-post", journalEntryController.CreateReversedJournalEntry)
	}
}
