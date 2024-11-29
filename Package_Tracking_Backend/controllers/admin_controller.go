package controllers

import (
	"Package_Tracking_Backend/db"
	"Package_Tracking_Backend/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterAdmin handles the registration of a new admin user
func RegisterAdmin(c *gin.Context) {

	// Verify if the requester has a super admin role
	requesterRole, exists := c.Get("role")
	if !exists || requesterRole != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var admin models.User

	// Parse the JSON input into the Admin struct
	if err := c.ShouldBindJSON(&admin); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Log the admin details without the password for security reasons
	log.Printf("Received admin registration data: Name=%s, Email=%s\n", admin.Name, admin.Email)

	admin.Role = "admin"
	// Hash the admin's password for secure storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register admin"})
		return
	}
	admin.Password = string(hashedPassword)

	// Prepare the SQL query to insert the new admin into the database
	query := `INSERT INTO users (name, email, password, role) 
              VALUES ($1, $2, $3, $4) RETURNING id`

	// Execute the query and get the new admin's ID
	var adminID int64
	err = db.DB.QueryRow(query, admin.Name, admin.Email, admin.Password, admin.Role).Scan(&adminID)
	if err != nil {
		log.Println("Error inserting admin into database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register admin"})
		return
	}

	// Respond with a success message and the admin's ID
	c.JSON(http.StatusOK, gin.H{"message": "Admin registered successfully", "admin_id": adminID})
}

func LoginAdmin(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Parse the JSON input into the loginData struct
	if err := c.ShouldBindJSON(&loginData); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Retrieve the admin's data from the database
	var admin models.Admin
	query := `SELECT id, name, email, password FROM admins WHERE email = $1`
	err := db.DB.Get(&admin, query, loginData.Email)
	if err != nil {
		log.Println("Admin not found:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginData.Password)); err != nil {
		log.Println("Password mismatch:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "admin_id": admin.ID, "name": admin.Name})
}

// GetAllOrders handles retrieving all orders for admin
func GetAllOrders(c *gin.Context) {
	var orders []models.Order
	query := `SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status 
              FROM orders`

	// Fetch all orders from the database
	err := db.DB.Select(&orders, query)
	if err != nil {
		log.Println("Error fetching orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateOrder handles updating the status of an existing order
func UpdateOrder(c *gin.Context) {
	var order models.Order

	// Bind the JSON request to the Order model
	if err := c.ShouldBindJSON(&order); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	orderID := c.Param("id")

	// Update the order in the database
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := db.DB.Exec(query, order.Status, orderID)
	if err != nil {
		log.Println("Error updating order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

// DeleteOrder deletes an order from the database
func DeleteOrder(c *gin.Context) {
	orderID := c.Param("id")

	// Delete the order from the database
	query := `DELETE FROM orders WHERE id = $1`

	_, err := db.DB.Exec(query, orderID)
	if err != nil {
		log.Println("Error deleting order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
