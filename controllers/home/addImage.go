package controllers

import (
	"ayana/db"
	"ayana/models"
	uploadClaudinary "ayana/utils/cloudinary-folder"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddImage(c *gin.Context) {

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

	newImageData := models.Image{
		Title: c.Request.PostFormValue("title"),
	}

	newImageData.Image = imageUrl

	result := db.DB.Debug().Create(&newImageData)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   newImageData,
	})

}
