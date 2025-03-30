package tests

import (
	"bytes"
	"file-sharing-backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router)
	return router
}

func TestUserRegistration(t *testing.T) {
	router := setupRouter()

	userData := map[string]string{
		"password": "dhairya",
		"email":    "dhairya@gmail.com",
	}

	jsonData, _ := json.Marshal(userData)

	req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}

func TestUserLogin(t *testing.T) {
	router := setupRouter()

	userData := map[string]string{
		"email": "dhairya@gmail.com",
		"password": "dhairya",
	}

	jsonData, _ := json.Marshal(userData)

	// Create a new POST request to login
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "access_token")
}

// TestAuthenticatedRoute tests a secured endpoint that requires JWT authentication
func TestAuthenticatedRoute(t *testing.T) {
	// First, we login to get the token
	loginData := map[string]string{
		"email": "dhairya@gmail.com",
		"password": "dhairya",
	}

	// Convert loginData to JSON
	loginJSON, _ := json.Marshal(loginData)

	// Setup the Gin router
	router := setupRouter()

	// Create a POST request to login
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert login is successful and contains the access token
	assert.Equal(t, http.StatusOK, w.Code)
	var loginResponse map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["access_token"]

	// Now, we test a secured route using the JWT token
	// Create a new GET request to access the protected route
	reqSecured := httptest.NewRequest(http.MethodGet, "/auth", nil)
	reqSecured.Header.Set("Authorization", "Bearer "+token)

	// Record the response
	wSecured := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(wSecured, reqSecured)

	// Assert the response
	assert.Equal(t, http.StatusOK, wSecured.Code)
	assert.Contains(t, wSecured.Body.String(), "authenticated user data")
}

// TestInvalidLogin tests invalid login credentials
func TestInvalidLogin(t *testing.T) {
	// Setup the Gin router
	router := setupRouter()

	// Prepare invalid user login data
	invalidData := map[string]string{
		"username": "invaliduser",
		"password": "wrongpassword",
	}

	// Convert invalidData to JSON
	jsonData, _ := json.Marshal(invalidData)

	// Create a new POST request to login
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")
}
