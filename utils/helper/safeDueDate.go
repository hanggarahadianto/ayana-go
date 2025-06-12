package helper

import "time"

// SafeDueDate untuk handle null due_date
func SafeDueDate(dueDate *time.Time) *time.Time {
	if dueDate != nil {
		return dueDate
	}
	return nil
}

func SafeRepaymentDate(repaymentDate *time.Time) *time.Time {
	if repaymentDate != nil {
		return repaymentDate
	}
	return nil
}
