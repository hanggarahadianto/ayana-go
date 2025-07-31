package handler

import (
	"ayana/dto"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// Struktur request payload
type DownloadJournalPayload struct {
	Entries          []dto.JournalEntryResponse `json:"entries"`
	Title            string                     `json:"title"`
	StartDate        string                     `json:"startDate"`
	EndDate          string                     `json:"endDate"`
	SelectedCategory string                     `json:"selectedCategory"`
	SearchTerm       string                     `json:"searchTerm"`
}

// Handler utama
func DownloadJournalTransactionHandler(c *gin.Context) {
	var payload DownloadJournalPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": "Payload tidak valid"})
		return
	}

	if len(payload.Entries) == 0 {
		c.JSON(400, gin.H{"error": "Data jurnal kosong"})
		return
	}

	file, err := generateExcelFile(payload)
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal membuat file Excel"})
		return
	}

	tmpFile, err := os.CreateTemp("", "journal-*.xlsx")
	if err != nil {
		c.JSON(500, gin.H{"error": "Gagal membuat file sementara"})
		return
	}
	defer os.Remove(tmpFile.Name())

	if err := file.SaveAs(tmpFile.Name()); err != nil {
		c.JSON(500, gin.H{"error": "Gagal menyimpan file Excel"})
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", `attachment; filename="journal-transaction.xlsx"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(tmpFile.Name())
}

// Membuat file Excel
func generateExcelFile(payload DownloadJournalPayload) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Jurnal"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Judul (baris 1-2)
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  16,
			Color: "#000000",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	titleText := strings.ToUpper(payload.Title)
	f.SetCellValue(sheetName, "A1", titleText)
	_ = f.MergeCell(sheetName, "A1", "K2")
	_ = f.SetCellStyle(sheetName, "A1", "K2", titleStyle)

	// Filter (baris 3)
	filterText := buildFilterText(payload)
	if filterText != "" {
		filterStyle, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: false,
				Size: 11,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "left",
				Vertical:   "center",
			},
		})
		f.SetCellValue(sheetName, "A4", filterText)
		_ = f.MergeCell(sheetName, "A4", "K4")
		_ = f.SetCellStyle(sheetName, "A4", "K4", filterStyle)
	}

	// Header (baris 6)
	headers := []string{
		"No", "Tanggal Transaksi", "Transaction ID", "Invoice", "Partner", "Kategori",
		"Nominal", "Deskripsi", "Tanggal Jatuh Tempo", "Tanggal Pelunasan", "Status Pembayaran",
	}
	for i, h := range headers {
		cell := fmt.Sprintf("%s6", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, h)
	}

	// Data (mulai dari baris 7)
	startRow := 7
	for i, j := range payload.Entries {
		row := startRow + i

		kategori := "-"
		if j.Status == "going" && j.DebitCategory != "" {
			kategori = j.DebitCategory
		} else if j.CreditCategory != "" {
			kategori = j.CreditCategory
		}

		status := ""
		switch j.Status {
		case "paid", "done":
			status = "Sudah Lunas"
		case "unpaid":
			status = "Belum Lunas"
		default:
			status = j.Status
		}

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), formatDateIndo(j.DateInputed))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), j.TransactionID)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), j.Invoice)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), j.Partner)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), kategori)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), j.Amount)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), j.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), formatDateIndo(j.DueDate))
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), formatDateIndo(j.RepaymentDate))
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), status)
	}

	// Border dan styling
	totalRows := len(payload.Entries) + 6 // baris 6 adalah header
	applyBorders(f, sheetName, "A6", fmt.Sprintf("K%d", totalRows))
	setColumnWidths(f, sheetName)

	return f, nil
}

// Format tanggal Bahasa Indonesia
var bulan = [...]string{
	"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

func formatDateIndo(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%02d %s %d", t.Day(), bulan[int(t.Month())], t.Year())
}

// Lebar kolom
func setColumnWidths(f *excelize.File, sheet string) {
	widths := map[string]float64{
		"A": 6,  // No
		"B": 18, // Tanggal Transaksi
		"C": 25, // Transaction ID
		"D": 20, // Invoice
		"E": 30, // Partner
		"F": 20, // Kategori
		"G": 15, // Nominal
		"H": 35, // Deskripsi
		"I": 18, // Tanggal Jatuh Tempo
		"J": 18, // Tanggal Pelunasan
		"K": 20, // Status
	}
	for col, width := range widths {
		_ = f.SetColWidth(sheet, col, col, width)
	}
}

// Tambahkan border ke range sel
func applyBorders(f *excelize.File, sheet, startCell, endCell string) {
	style, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err == nil {
		_ = f.SetCellStyle(sheet, startCell, endCell, style)
	}
}

// Filter text dinamis
func buildFilterText(p DownloadJournalPayload) string {
	filter := ""

	if p.StartDate != "" && p.EndDate != "" {
		filter += fmt.Sprintf("Periode: %s - %s", p.StartDate, p.EndDate)
	}

	if p.SelectedCategory != "" {
		if filter != "" {
			filter += " | "
		}
		filter += "Kategori: " + p.SelectedCategory
	}

	if p.SearchTerm != "" {
		if filter != "" {
			filter += " | "
		}
		filter += "Pencarian: " + p.SearchTerm
	}

	return filter
}
