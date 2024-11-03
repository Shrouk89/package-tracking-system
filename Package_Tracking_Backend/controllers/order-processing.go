package controllers

import (
	"Package_Tracking_Backend/db"
	"Package_Tracking_Backend/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateOrder handles order creation
func CreateOrder(c *gin.Context) {
	var order models.Order

	// Bind JSON request to the order model
	if err := c.ShouldBindJSON(&order); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Insert order into the database
	query := `INSERT INTO orders (user_id, pickup_location, dropoff_location, package_details, delivery_time, status)
              VALUES ($1, $2, $3, $4, $5, 'pending') RETURNING id`

	var orderID int64
	err := db.DB.QueryRow(query, order.UserID, order.PickupLocation, order.DropoffLocation, order.PackageDetails, order.DeliveryTime).Scan(&orderID)
	if err != nil {
		log.Println("Error inserting order into database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order_id": orderID})
}


// GetOrdersByUser handles retrieving all orders for a user
func GetOrdersByUser(c *gin.Context) {
	userID := c.Query("user_id") // Get user ID from query parameters

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var orders []models.Order
	query := `SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status 
              FROM orders WHERE user_id = $1`
	err := db.DB.Select(&orders, query, userID)

	if err != nil {
		log.Println("Error fetching orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}
