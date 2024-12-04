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

	// Retrieve user_id from the context (set by authentication middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		log.Println("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Bind JSON request to the order model
	if err := c.ShouldBindJSON(&order); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Assign the authenticated userID to the order
	order.UserID = userID.(int64)

	// Insert the new order into the database
	query := `INSERT INTO orders (user_id, pickup_location, dropoff_location, package_details, delivery_time, status)
              VALUES ($1, $2, $3, $4, $5, 'pending') RETURNING id`

	var orderID int64
	err := db.DB.QueryRow(query, order.UserID, order.PickupLocation, order.DropoffLocation, order.PackageDetails, order.DeliveryTime).Scan(&orderID)
	if err != nil {
		log.Println("Error inserting order into database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Return success response with order ID
	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order_id": orderID})
}

// GetOrdersByUser handles retrieving all orders for an authenticated user
func GetOrdersByUser(c *gin.Context) {
	// Retrieve user_id from the context (set by authentication middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		log.Println("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Fetch orders from the database for the authenticated user
	var orders []models.Order
	query := `SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status 
              FROM orders WHERE user_id = $1`
	err := db.DB.Select(&orders, query, userID.(int64))

	if err != nil {
		log.Println("Error fetching orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	// Return the retrieved orders
	c.JSON(http.StatusOK, orders)
}

// GetOrderDetails handles retrieving the details of a specific order by its ID
func GetOrderDetails(c *gin.Context) {
	orderID := c.Param("id") // Get order ID from URL parameters

	// Fetch the specific order details from the database
	var order models.Order
	query := `SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status 
              FROM orders WHERE id = $1`
	err := db.DB.Get(&order, query, orderID)

	if err != nil {
		log.Println("Error fetching order details:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order details"})
		return
	}

	// Return the order details
	c.JSON(http.StatusOK, order)
}

// CancelOrder handles the cancellation of a pending order
func CancelOrder(c *gin.Context) {
	orderID := c.Param("id") // Get order ID from URL parameters

	// Check if the order is pending before allowing cancellation
	var status string
	var order models.Order
	queryCheck := `SELECT status FROM orders WHERE id = $1`
	err := db.DB.Get(&order, queryCheck, orderID)
	if err != nil {
		log.Println("Error checking order status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check order status"})
		return
	}

	if status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only pending orders can be cancelled"})
		return
	}

	// Update order status to cancelled
	queryUpdate := `UPDATE orders SET status = 'cancelled' WHERE id = $1`
	_, err = db.DB.Exec(queryUpdate, orderID)
	if err != nil {
		log.Println("Error updating order status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	// Return success response after cancellation
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}
