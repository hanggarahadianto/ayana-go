package controllers

// import (
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"ayana/db"
// 	"ayana/models"
// 	uploadClaudinary "ayana/utils/cloudinary-folder"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// func CreateHome(c *gin.Context) {
// 	// Validate required fields
// 	requiredFields := []string{"title", "location", "content", "address", "bathroom", "bedroom", "square", "price", "quantity", "status"}
// 	missingFields := []string{}

// 	for _, field := range requiredFields {
// 		if c.Request.PostFormValue(field) == "" {
// 			missingFields = append(missingFields, field)
// 		}
// 	}

// 	if len(missingFields) > 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  "failed",
// 			"message": "Missing required fields: " + strings.Join(missingFields, ", "),
// 		})
// 		return
// 	}

// 	// File upload validation
// 	file, header, err := c.Request.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status": "failed",
// 			"error":  "file is required",
// 		})
// 		return
// 	}
// 	defer file.Close()

// 	filename := header.Filename
// 	var newHome models.Home

// 	// Mulai transaksi database
// 	err = db.DB.Transaction(func(tx *gorm.DB) error {
// 		imageUrl, err := uploadClaudinary.UploadToCloudinary(file, filename)
// 		if err != nil {
// 			return err // Jika gagal upload, transaksi otomatis rollback
// 		}

// 		// Convert price, quantity, sequence safely
// 		price, err := strconv.Atoi(c.Request.PostFormValue("price"))
// 		if err != nil {
// 			_ = uploadClaudinary.DeleteFromCloudinary(imageUrl) // Hapus gambar jika terjadi error
// 			return err
// 		}

// 		quantity, err := strconv.Atoi(c.Request.PostFormValue("quantity"))
// 		if err != nil {
// 			_ = uploadClaudinary.DeleteFromCloudinary(imageUrl)
// 			return err
// 		}

// 		sequence, err := strconv.Atoi(c.Request.PostFormValue("sequence"))
// 		if err != nil {
// 			_ = uploadClaudinary.DeleteFromCloudinary(imageUrl)
// 			return err
// 		}

// 		// Save to database
// 		now := time.Now()
// 		newHome = models.Home{
// 			Title:     c.Request.PostFormValue("title"),
// 			Location:  c.Request.PostFormValue("location"),
// 			Content:   c.Request.PostFormValue("content"),
// 			Address:   c.Request.PostFormValue("address"),
// 			Bathroom:  c.Request.PostFormValue("bathroom"),
// 			Bedroom:   c.Request.PostFormValue("bedroom"),
// 			Square:    c.Request.PostFormValue("square"),
// 			Price:     float64(price),
// 			Quantity:  quantity,
// 			Status:    c.Request.PostFormValue("status"),
// 			Sequence:  sequence,
// 			CreatedAt: now,
// 			UpdatedAt: now,
// 			Image:     imageUrl,
// 		}

// 		// Insert into database
// 		if err := tx.Create(&newHome).Error; err != nil {
// 			_ = uploadClaudinary.DeleteFromCloudinary(imageUrl) // Hapus gambar jika gagal menyimpan ke DB
// 			return err
// 		}

// 		return nil // Jika berhasil, commit transaksi
// 	})

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  "failed",
// 			"message": err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  "success",
// 		"message": "Home created successfully",
// 		"data":    newHome,
// 	})
// }
