package journalentry

// import (
// 	"ayana/db"
// 	"ayana/models"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/google/uuid"
// )

// func CreateJournalEntry(c *gin.Context) {
// 	var journal models.JournalEntry

// 	// Parse JSON dari request
// 	if err := c.ShouldBindJSON(&journal); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal parsing JSON: " + err.Error()})
// 		return
// 	}

// 	var existing models.JournalEntry
// 	if err := db.DB.Where("invoice = ?", journal.Invoice).First(&existing).Error; err == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invoice tidak boleh sama"})
// 		return
// 	}

// 	// Generate UUID jika belum ada
// 	if journal.ID == uuid.Nil {
// 		journal.ID = uuid.New()
// 	}

// 	// Set tanggal jika null (optional)
// 	if journal.Date == nil {
// 		now := time.Now()
// 		journal.Date = &now
// 	}

// 	// Set waktu dibuat dan diupdate
// 	journal.CreatedAt = time.Now()
// 	journal.UpdatedAt = time.Now()

// 	// Validasi: minimal ada 1 baris jurnal
// 	if len(journal.Lines) == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimal satu baris jurnal (lines) diperlukan."})
// 		return
// 	}

// 	// Set ID dan JournalID untuk tiap JournalLine
// 	for i := range journal.Lines {
// 		if journal.Lines[i].ID == uuid.Nil {
// 			journal.Lines[i].ID = uuid.New()
// 		}
// 		journal.Lines[i].JournalID = journal.ID
// 		journal.Lines[i].CreatedAt = time.Now()
// 		journal.Lines[i].UpdatedAt = time.Now()
// 	}

// 	// Simpan ke database (akan insert ke JournalEntry dan JournalLine)
// 	if err := db.DB.Create(&journal).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jurnal: " + err.Error()})
// 		return
// 	}

// 	// Respons sukses
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  "success",
// 		"message": "Journal entry berhasil dibuat",
// 		"data":    journal,
// 	})
// }
