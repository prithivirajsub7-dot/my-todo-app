// Token verify
package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization") //JWT tokens inga vanthuruku
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"}) //unauthorization ah nu check pannum
			c.Abort()//unauthorization na middleware la irunthu request stop pannum
			return
		}

		// Format: Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1] //JWT token ah extract panro

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { //jwt.Parse valiya token ah verify panro
			return jwtSecret, nil //jwtSecret →  secret key , token ah sign pannumpodhu use pannathu.
		})

		if err != nil || !token.Valid { //token invalidnu or error vantha 401 kudukum
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
        claims := token.Claims.(jwt.MapClaims)
        userIDFloat := claims["user_id"].(float64)
        c.Set("userID", uint(userIDFloat)) // uint-க்கு safe conversion
		

		c.Next()//Token validna, indha request ah next handlerku send panrom
	}
}