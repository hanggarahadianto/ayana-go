package helper

import (
	"fmt"
	"time"
)

func HitungPaymentNote(assetOrDebtStatus string, due *time.Time, repayment *time.Time, now time.Time) (string, string) {
	if due == nil || due.IsZero() {
		return "", ""
	}

	// Untuk debt yang sudah dibayar
	if assetOrDebtStatus == "done" && repayment != nil && !repayment.IsZero() {
		diff := due.Sub(*repayment).Hours() / 24
		if diff >= 0 {
			return fmt.Sprintf("Dibayar %d Hari Sebelum Jatuh Tempo", int(diff)), "teal"
		}
		return fmt.Sprintf("Dibayar %d Hari Setelah Jatuh Tempo", int(-diff)), "red"
	}

	// Untuk debt yang masih berjalan
	if assetOrDebtStatus == "going" {
		diff := due.Sub(now).Hours() / 24
		if diff >= 0 {
			return fmt.Sprintf("Pembayaran Kurang %d Hari Lagi", int(diff)), "blue"
		}
		return fmt.Sprintf("Pembayaran Terlambat %d Hari Yang Lalu", int(-diff)), "red"
	}

	return "", ""
}
