package helper

import (
	"ayana/dto"
	lib "ayana/lib"
	"time"
)

func EnrichJournalEntryResponses(responses []dto.JournalEntryResponse, entryType string, now time.Time) []dto.JournalEntryResponse {
	for i, line := range responses {
		note, color := lib.HitungPaymentNote(line.DueDate, line.RepaymentDate, entryType, now)
		responses[i].PaymentNote = note
		responses[i].PaymentNoteColor = color
	}
	return responses
}
