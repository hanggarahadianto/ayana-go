package helper

import (
	"ayana/dto"
	lib "ayana/lib"
	"time"
)

func EnrichJournalEntryResponses(responses []dto.JournalEntryResponse, status string, now time.Time) []dto.JournalEntryResponse {
	for i, line := range responses {
		note, color := lib.HitungPaymentNote(status, line.DueDate, line.RepaymentDate, now)
		responses[i].PaymentNote = note
		responses[i].PaymentNoteColor = color
	}
	return responses
}
