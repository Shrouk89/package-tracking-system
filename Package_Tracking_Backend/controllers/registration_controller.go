package controllers

import (
	"Package_Tracking_Backend/db"
	"Package_Tracking_Backend/models"
	"Package_Tracking_Backend/utils2"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration
func Register(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Default role to 'user' unless explicitly set to 'admin' (backend-controlled)
	if user.Role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot self-register as an admin."})
		return
	}

	// Insert the new user into the database
	query := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id`
	var userID int64
	err = db.DB.QueryRow(query, user.Name, user.Email, string(hashedPassword), user.Role).Scan(&userID)
	if err != nil {
		log.Println("Error inserting user into the database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful",
		"user_id": userID,
	})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Retrieve the user from the database
	var user models.User
	query := `SELECT id, name, email, password, role FROM users WHERE email = $1`
	err := db.DB.Get(&user, query, credentials.Email)
	if err != nil {
		log.Println("Error fetching user:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// // Check if user is an admin
	// if user.Role != "admin" {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	// 	return
	// }

	// Generate JWT token
	token, err := utils2.GenerateJWT(user.ID, user.Role)
	if err != nil {
		log.Println("Error generating JWT:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"role":    user.Role,
		"user_id": user.ID,
	})
}

func ApproveAdmin(c *gin.Context) {
	// Verify super admin role
	requesterRole, exists := c.Get("role")
	if !exists || requesterRole != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get the user ID from the URL
	userID := c.Param("id")

	// Update the role from 'pending_admin' to 'admin'
	query := `UPDATE users SET role = 'admin' WHERE id = $1 AND role = 'pending_admin'`
	_, err := db.DB.Exec(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin approved successfully"})
}
