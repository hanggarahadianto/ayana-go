package helper

import "ayana/models"

func IsValidStatus(s string) bool {
	switch models.Status(s) {
	case models.StatusDraft, models.StatusApproved, models.StatusPaid, models.StatusUnpaid, models.StatusCancelled:
		return true
	default:
		return false
	}
}

func IsValidTransactionType(t string) bool {
	switch models.TransactionType(t) {
	case models.PayinType, models.PayoutType:
		return true
	default:
		return false
	}
}
