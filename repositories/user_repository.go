package repositories

import (
	"database/sql"
	"errors"
	"log"

	"file-sharing-backend/config"
	"file-sharing-backend/models"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser inserts a new user into the database
func CreateUser(email, password string) error {
	// Check if DB is initialized
	if config.DB == nil {
		log.Println("ERROR: Database connection is nil!")
		return errors.New("database not connected")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("ERROR: Failed to hash password:", err)
		return err
	}

	// Insert into DB
	_, err = config.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, hashedPassword)
	if err != nil {
		log.Println("ERROR: Failed to insert user into database:", err)
		return err
	}

	log.Println("SUCCESS: User registered:", email)
	return nil
}

// GetUserByEmail fetches a user by email
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := config.DB.QueryRow("SELECT id, email, password FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		log.Println("WARNING: User not found:", email)
		return nil, errors.New("user not found")
	} else if err != nil {
		log.Println("ERROR: Query failed:", err)
		return nil, err
	}

	log.Println("SUCCESS: User retrieved:", email)
	return &user, nil
}
