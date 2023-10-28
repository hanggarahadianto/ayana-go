package middlewares

// import (
// 	"ayana/db"
// 	"ayana/utils"
// 	"fmt"

// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// )

// func DeserializeUser() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var token string
// 		cookie, err := c.Cookie("token")

// 		authorizationHeader := c.Request.Header.Get("Authorization")
// 		fields := strings.Fields(authorizationHeader)

// 		if len(fields) != 0 && fields[0] == "Bearer" {
// 			token = fields[1]
// 		} else if err == nil {
// 			token = cookie
// 		}

// 		if token == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"status":  "fail",
// 				"message": "You are not logged in",
// 			})
// 			return
// 		}

// 		config, _ := utils.LoadConfig(".")
// 		sub, err := ValidateToken(token, config.AccessTokenPublicKey)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"status":  "fail",
// 				"message": err.Error(),
// 			})
// 			return
// 		}

// 		var user models.User
// 		result := db.DB.First(&user, "id = ?", fmt.Sprint(sub))
// 		if result.Error != nil {
// 			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
// 				"status":  "fail",
// 				"message": "the user belonging to this token no logger exists",
// 			})
// 		}

// 		c.Set("currentUser", user)
// 		c.Next()

// 	}
// }
