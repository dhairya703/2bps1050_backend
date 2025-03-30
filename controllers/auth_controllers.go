package controllers

import (
	"file-sharing-backend/repositories"
	"file-sharing-backend/models"
	"file-sharing-backend/config"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
   "os"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))


// Register handles user registration
func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call repository function to create user
	err := repositories.CreateUser(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Generate JWT Token
func generateToken(userID int) (string, error) {
	// Set claims
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 Days Expiry
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch user from database
	user, err := repositories.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT Token
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Insert session into database
	query := "INSERT INTO sessions (user_id, token, created_at, expires_at) VALUES ($1, $2, NOW(), $3)"
	_, err = config.DB.Exec(query, user.ID, token, time.Now().Add(7*24*time.Hour))

	// Debugging output
	if err != nil {
	
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Could not store session",
			"details": err.Error(), // Show exact database error
		})
		return
	}

	// Successful login
	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Login successful",
	})
}

func Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Protected route accessed!", "user_id": userID})
}
