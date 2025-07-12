package handler

import (
	"ayana/db"
	"ayana/models"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// UploadPresenceHandler handles Excel presence upload
func UploadPresenceHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File tidak ditemukan"})
		return
	}

	tmpFilePath, err := saveTempFile(file)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal menyimpan file"})
		return
	}
	defer os.Remove(tmpFilePath)

	f, err := excelize.OpenFile(tmpFilePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal membuka file Excel"})
		return
	}

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		c.JSON(500, gin.H{"error": "Sheet tidak ditemukan dalam Excel"})
		return
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal membaca isi sheet"})
		return
	}

	total, sukses, gagal := processPresenceRows(rows)
	message := fmt.Sprintf("Import selesai: total=%d, sukses=%d, gagal=%d", total, sukses, gagal)
	c.JSON(200, gin.H{"message": message})
}

// saveTempFile saves the uploaded file to a temporary path
func saveTempFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	tmpFile, err := os.CreateTemp("", "presence-*.xlsx")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, src)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func processPresenceRows(rows [][]string) (total int, sukses int, gagal int) {
	for i, row := range rows {
		total++

		if i == 0 {
			continue // Skip header
		}

		if len(row) < 4 {
			fmt.Printf("❌ Baris %d dilewati: kolom tidak lengkap (%d kolom)\n", i+1, len(row))
			gagal++
			continue
		}

		tanggalStr := strings.TrimSpace(row[0])
		rawTanggal := strings.TrimSpace(row[1])
		jam := strings.TrimSpace(row[2])
		nama := strings.TrimSpace(row[3])

		if nama == "" {
			fmt.Printf("❌ Baris %d dilewati: nama kosong\n", i+1)
			gagal++
			continue
		}

		// Debug info
		fmt.Printf("🔍 Baris %d | Nama: %s | Tanggal: %s | Jam: %s\n", i+1, nama, tanggalStr, jam)

		employee, err := findEmployeeByName(nama)
		if err != nil {
			fmt.Printf("❌ Baris %d | Gagal cari karyawan: %s\n", i+1, err.Error())
			gagal++
			continue
		}

		tanggalScan, err := parseTanggal(tanggalStr)
		if err != nil {
			fmt.Printf("❌ Baris %d | Format tanggal salah '%s': %v\n", i+1, tanggalStr, err)
			gagal++
			continue
		}

		dayName := strings.ToLower(tanggalScan.Weekday().String())

		presence := models.Presence{
			EmployeeID: employee.ID,
			CompanyID:  employee.CompanyID,
			ScanDate:   tanggalScan,
			ScanTime:   jam,
			RawDate:    rawTanggal,
			Day:        dayName,
		}

		if err := db.DB.Create(&presence).Error; err != nil {
			fmt.Printf("❌ Baris %d | Gagal simpan presensi: %v\n", i+1, err)
			gagal++
			continue
		}

		fmt.Printf("✅ Baris %d | Presensi berhasil disimpan untuk %s\n", i+1, employee.Name)
		sukses++
	}

	return
}

func findEmployeeByName(name string) (*models.Employee, error) {
	var emp models.Employee

	// Buat pattern pencarian: %rafli%
	keyword := "%" + strings.ToLower(name) + "%"

	err := db.DB.
		Where("LOWER(name) ILIKE ?", keyword).
		First(&emp).Error

	if err != nil {
		return nil, fmt.Errorf("tidak ditemukan nama mirip '%s'", name)
	}

	return &emp, nil
}

// parseTanggal converts string to time.Time using expected format
func parseTanggal(str string) (time.Time, error) {
	str = strings.TrimSpace(str)
	formats := []string{
		"2006-01-02",          // 2025-07-09
		"02-01-2006",          // 09-07-2025
		"02/01/2006",          // 09/07/2025
		"2006/01/02",          // 2025/07/09
		"02-01-2006 15:04:05", // 09-07-2025 08:35:49
		"02/01/2006 15:04:05", // 09/07/2025 08:35:49
		"2006-01-02 15:04:05", // 2025-07-09 08:35:49
		"2006/01/02 15:04:05", // 2025/07/09 08:35:49
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, str)
		if err == nil {
			return t, nil
		}
		lastErr = err
	}

	return time.Time{}, fmt.Errorf("format tanggal tidak dikenali: '%s', error terakhir: %v", str, lastErr)
}
