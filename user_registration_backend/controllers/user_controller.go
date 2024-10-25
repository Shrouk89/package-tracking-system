package controllers

import (
	"log"
	"net/http"
	"user_registration_backend/db"
	"user_registration_backend/models"

	"github.com/gin-gonic/gin"
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
