package controllers

import (
	"mime/multipart"
	"net/http"

	"ayana/db"
	"ayana/models"
	uploadClaudinary "ayana/utils/cloudinary-folder"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateHome(c *gin.Context) {

	filename, ok := c.Get("filePath")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "filename not found",
			"data":  filename,
		})
	}
	file, ok := c.Get("file")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "failed",
		})
	}

	// upload file
	imageUrl, err := uploadClaudinary.UploadtoHomeFolder(file.(multipart.File), filename.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "failed",
			"error":  err.Error()})
		return
	}

	now := time.Now()
	newHome := models.Home{
		Title:    c.Request.PostFormValue("title"),
		Content:  c.Request.PostFormValue("content"),
		Address:  c.Request.PostFormValue("address"),
		Bathroom: c.Request.PostFormValue("bathroom"),
		Bedroom:  c.Request.PostFormValue("bedroom"),
		Square:   c.Request.PostFormValue("square"),

		CreatedAt: now,
		UpdatedAt: now,
	}

	newHome.Image = imageUrl

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
