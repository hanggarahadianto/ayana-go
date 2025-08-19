package service

import (
	"ayana/models"
	"fmt"
	"strconv"
	"time"
)

func EnrichProjectStatus(p models.Project) models.ProjectWithStatus {
	// Kalau sudah selesai
	if p.ProjectStatus == "done" && p.ProjectStart != nil && p.ProjectFinished != nil {
		durasi := int(p.ProjectFinished.Sub(*p.ProjectStart).Hours() / 24) // pakai ProjectFinished
		targetDurasi := 0
		if p.ProjectTime != "" {
			if d, err := strconv.Atoi(p.ProjectTime); err == nil {
				targetDurasi = d
			}
		}

		delay := durasi - targetDurasi
		finishStatus := "Tepat waktu"

		if delay > 0 {
			finishStatus = fmt.Sprintf("Terlambat %d hari", delay)
		} else if delay < 0 {
			finishStatus = fmt.Sprintf("Lebih cepat %d hari", -delay)
		}

		isOnTime := delay <= 0 // true kalau tepat waktu ATAU lebih cepat

		color := "red"
		if isOnTime {
			color = "blue"
		}

		return models.ProjectWithStatus{
			Project:      p,
			StatusText:   fmt.Sprintf("Selesai dalam %d hari", durasi),
			Color:        color,
			IsOnTime:     &isOnTime,
			DelayDays:    delay,
			FinishStatus: finishStatus,
		}
	}

	// Kalau masih berjalan
	if p.ProjectStart != nil {
		now := time.Now().UTC()
		durasiBerjalan := int(now.Sub(*p.ProjectStart).Hours() / 24)

		return models.ProjectWithStatus{
			Project:    p,
			StatusText: fmt.Sprintf("Sedang berjalan %d hari", durasiBerjalan),
			Color:      "green",
		}
	}

	// Kalau belum ada tanggal mulai
	return models.ProjectWithStatus{
		Project:    p,
		StatusText: "Belum dimulai",
		Color:      "yellow",
	}
}
