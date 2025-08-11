package middlewares

import (
	utilsEnv "ayana/utils/env"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil token dari cookie
		tokenStr, err := c.Cookie("token")
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token cookie required"})
			return
		}

		config, _ := utilsEnv.LoadConfig(".")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			fmt.Printf("DEBUG - token claims: %#v\n", claims)
			if username, ok := claims["username"].(string); ok {
				fmt.Println("DEBUG - username dari token:", username)
				c.Set("username", username)
			} else {
				fmt.Println("DEBUG - username di token tidak ditemukan atau bukan string")
			}
		}

		c.Next()
	}
}
