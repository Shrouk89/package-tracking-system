package controllers

import (
	"Package_Tracking_Backend/db"
	"Package_Tracking_Backend/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
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

	// Log the user details without the password
	log.Printf("Received user registration data: Name=%s, Email=%s, Phone=%s\n", user.Name, user.Email, user.Phone)

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}
	user.Password = string(hashedPassword)

	// Insert user into the database
	query := `INSERT INTO users (name, email, phone, password) 
              VALUES ($1, $2, $3, $4) RETURNING id`

	var userID int64
	err = db.DB.QueryRow(query, user.Name, user.Email, user.Phone, user.Password).Scan(&userID)
	if err != nil {
		log.Println("Error inserting user into database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Return success response with user ID
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user_id": userID})
}

var jwtSecretKey = []byte("UzN4041qFo+9TKhXaNeAqP/DJ8btfIGIT1rWsO6CyC8=") // Use an environment variable for security

// LoginUser handles user login
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
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		log.Println("Password mismatch:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 1-day expiration
	})
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return user data along with token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
		"token": tokenString,
	})
}
