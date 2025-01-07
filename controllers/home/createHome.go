package controllers

import (
	"net/http"
	"strconv"

	"ayana/db"
	"ayana/models"
	uploadClaudinary "ayana/utils/cloudinary-folder"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateHome(c *gin.Context) {

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "file not found",
		})
		return
	}
	defer file.Close()

	// Retrieve the filename
	filename := header.Filename

	// upload file
	imageUrl, err := uploadClaudinary.UploadtoHomeFolder(file, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed on upload cloudinary",
			"error":  err.Error(),
		})
		return
	}

	price, err := strconv.Atoi(c.Request.PostFormValue("price"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "invalid price",
		})
		return
	}

	quantity, err := strconv.Atoi(c.Request.PostFormValue("quantity"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
			"error":  "invalid quantity",
		})
		return
	}

	status := models.StatusType(c.Request.PostFormValue("status"))

	now := time.Now()
	newHome := models.Home{
		Title:    c.Request.PostFormValue("title"),
		Content:  c.Request.PostFormValue("content"),
		Address:  c.Request.PostFormValue("address"),
		Bathroom: c.Request.PostFormValue("bathroom"),
		Bedroom:  c.Request.PostFormValue("bedroom"),
		Square:   c.Request.PostFormValue("square"),
		Price:    price,
		Quantity: quantity,
		Status:   status,

		CreatedAt: now,
		UpdatedAt: now,
	}

	newHome.Image = imageUrl

	db.DB.Exec("DISCARD ALL")

	result := db.DB.Debug().Create(&newHome)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newHome,
	})

}
