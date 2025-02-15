package controllers

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"ayana/db"
	"ayana/models"
	uploadClaudinary "ayana/utils/cloudinary-folder"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateHome(c *gin.Context) {
	// Validate required fields
	requiredFields := []string{"title", "location", "content", "address", "bathroom", "bedroom", "square", "price", "quantity", "status"}
	missingFields := []string{}

	for _, field := range requiredFields {
		if c.Request.PostFormValue(field) == "" {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Missing required fields: " + strings.Join(missingFields, ", "),
		})
		return
	}

	// File upload validation
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "file is required",
		})
		return
	}
	defer file.Close()

	title := c.Request.PostFormValue("title")
	slugTitle := strings.ReplaceAll(strings.ToLower(title), " ", "-") // Convert to slug
	uniqueID := uuid.New().String()[:8]                               // Short unique ID
	ext := filepath.Ext(header.Filename)                              // Get file extension
	filename := slugTitle + "-" + uniqueID + ext                      // Create final filename

	// Upload file

	imageUrl, err := uploadClaudinary.UploadtoHomeFolder(file, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  "upload to cloudinary failed: " + err.Error(),
		})
		return
	}

	// Convert price and quantity safely
	price, err := strconv.Atoi(c.Request.PostFormValue("price"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "price must be a valid number",
		})
		return
	}

	quantity, err := strconv.Atoi(c.Request.PostFormValue("quantity"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "quantity must be a valid number",
		})
		return
	}

	// Save to database
	now := time.Now()
	newHome := models.Home{
		Title:     c.Request.PostFormValue("title"),
		Location:  c.Request.PostFormValue("location"),
		Content:   c.Request.PostFormValue("content"),
		Address:   c.Request.PostFormValue("address"),
		Bathroom:  c.Request.PostFormValue("bathroom"),
		Bedroom:   c.Request.PostFormValue("bedroom"),
		Square:    c.Request.PostFormValue("square"),
		Price:     float64(price),
		Quantity:  quantity,
		Status:    c.Request.PostFormValue("status"),
		CreatedAt: now,
		UpdatedAt: now,
		Image:     imageUrl,
	}

	// Insert into database
	result := db.DB.Debug().Create(&newHome)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": result.Error.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newHome,
	})
}
