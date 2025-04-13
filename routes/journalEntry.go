package routes

import (
	journalEntryController "ayana/controllers/journalEntry"
	// "ayana/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupJournalEntryRouter(r *gin.Engine) {
	journalEntry := r.Group("/journal-entry")
	{
		journalEntry.POST("/post", journalEntryController.CreateJournalEntry)

	}
}
