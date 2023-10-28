package controllers

import (
	"ayana/db"
	"ayana/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateHome(c *gin.Context) {
	homeId := c.Param("id")

	var updateHomeData models.Home
	result := db.DB.Debug().First(&updateHomeData, "id = ?", homeId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "failed",
			"message": "home id not found",
		})
		return
	}

	// filename, ok := c.Get("filePath")
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "filename not found",
	// 		"data":  filename,
	// 	})
	// }
	// file, ok := c.Get("file")
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status": "failed",
	// 	})
	// }

	// upload file
	// imageUrl, err := uploadClaudinary.UploadtoHomeFolder(file.(multipart.File), filename.(string))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status": "failed",
	// 		"error":  err.Error()})
	// 	return
	// }

	afterUpdateHome := models.Home{
		ID:       updateHomeData.ID,
		Title:    c.Request.PostFormValue("title"),
		Content:  c.Request.PostFormValue("content"),
		Address:  c.Request.PostFormValue("address"),
		Bathroom: c.Request.PostFormValue("bathroom"),
		Bedroom:  c.Request.PostFormValue("bedroom"),
		Square:   c.Request.PostFormValue("square"),
	}

	// newHome.Image = imageUrl

	db.DB.Debug().Model(&updateHomeData).Updates(afterUpdateHome)
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  "success",
		"message": afterUpdateHome,
	})

}
