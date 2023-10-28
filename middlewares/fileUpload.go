package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FileUploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Bad request",
			})
			return
		}
		defer file.Close()

		// pass the file and its name to the controller
		c.Set("filePath", header.Filename)
		c.Set("file", file)

		// continue to controller
		c.Next()
	}
}
