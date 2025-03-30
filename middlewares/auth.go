package middlewares

import (
	"file-sharing-backend/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("🔍 Received Auth Header:", authHeader)

		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Println("Missing or Invalid Token Prefix")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("Extracted Token:", tokenString)

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("Invalid JWT token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract user_id from token
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			log.Println("Error extracting user_id from token")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}
		userID := int(userIDFloat)
		fmt.Println("Extracted User ID:", userID)

		// Validate session from database
		var count int
		err = config.DB.QueryRow("SELECT COUNT(*) FROM sessions WHERE token = $1 AND expires_at > $2", tokenString, time.Now()).Scan(&count)

		if err != nil {
			log.Println("Database error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}

		if count == 0 {
			log.Println("Session not found or expired")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
			c.Abort()
			return
		}

		// Store user_id in context for later use
		c.Set("user_id", userID)
		fmt.Println("Authentication successful for user:", userID)

		c.Next()
	}
}
