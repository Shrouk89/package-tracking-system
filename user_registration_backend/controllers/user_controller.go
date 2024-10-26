package controllers

import (
	"log"
	"net/http"
	"user_registration_backend/db"
	"user_registration_backend/models"

	"github.com/gin-gonic/gin"
        "golang.org/x/crypto/bcrypt"  // ADDED: Import bcrypt for password hashing

)

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var user models.User

	// Bind JSON request to the user model
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Printf("Received user registration data: %+v\n", user)

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) // ADDED: Hash password before storing
	if err != nil {
		log.Println("Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	user.Password = string(hashedPassword)  // ADDED: Store hashed password


	// Insert user into the database
	query := `INSERT INTO users (name, email, phone, password) 
              VALUES ($1, $2, $3, $4) RETURNING id`

	var userID int64
	err := db.DB.QueryRow(query, user.Name, user.Email, user.Phone, user.Password).Scan(&userID)

	if err != nil {
		log.Println("Error inserting user into database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user_id": userID})
}

// LoginUser handles user login - ADDED: New function for handling user login
func LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Bind JSON request to login data
	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Retrieve user data from the database
	var user models.User
	query := `SELECT id, name, email, phone, password FROM users WHERE email = $1`
	err := db.DB.Get(&user, query, loginData.Email)

	if err != nil {
		log.Println("User not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil { // ADDED: Compare hashed password
		log.Println("Password mismatch:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": user.ID, "name": user.Name})
}
