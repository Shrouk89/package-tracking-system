package controllers

import (
	"log"
	"net/http"

	"Package_Tracking_Backend/db"
	"Package_Tracking_Backend/models"

	"github.com/gin-gonic/gin"
)

func GetPendingOrders(c *gin.Context) {
	var orders []models.Order
	query := `SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status
              FROM orders WHERE courier_id IS NULL AND status = 'pending'`
	err := db.DB.Select(&orders, query)
	if err != nil {
		log.Println("Error fetching pending orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func AcceptOrder(c *gin.Context) {
	courierID := c.GetInt64("user_id") // Extract courier ID from the context
	orderID := c.Param("id")           // Order ID from the URL

	// Check if the courier has already accepted 2 orders
	var acceptedCount int
	queryCheck := `SELECT COUNT(*) FROM orders WHERE courier_id = $1 AND status = 'accepted'`
	err := db.DB.Get(&acceptedCount, queryCheck, courierID)
	if err != nil {
		log.Println("Error checking courier's accepted orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check courier's orders"})
		return
	}

	if acceptedCount >= 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only accept up to 2 orders at a time"})
		return
	}

	// Update the order to assign it to the courier and mark it as accepted
	queryUpdate := `UPDATE orders SET courier_id = $1, status = 'accepted' WHERE id = $2`
	_, err = db.DB.Exec(queryUpdate, courierID, orderID)
	if err != nil {
		log.Println("Error accepting order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order accepted successfully"})
}

func DeclineOrder(c *gin.Context) {
	orderID := c.Param("id") // Order ID from the URL

	query := `UPDATE orders SET status = 'declined' WHERE id = $1`
	_, err := db.DB.Exec(query, orderID)
	if err != nil {
		log.Println("Error declining order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decline order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order declined successfully"})
}

func UpdateOrderStatus(c *gin.Context) {
	courierID := c.GetInt64("user_id") // Extract courier ID from the context
	orderID := c.Param("id")           // Order ID from the URL

	var requestBody struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Ensure the order belongs to the courier
	queryCheck := `SELECT COUNT(*) FROM orders WHERE id = $1 AND courier_id = $2`
	var count int
	err := db.DB.Get(&count, queryCheck, orderID, courierID)
	if err != nil || count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not assigned to this order"})
		return
	}

	// Update the order status
	queryUpdate := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err = db.DB.Exec(queryUpdate, requestBody.Status, orderID)
	if err != nil {
		log.Println("Error updating order status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// AddCourier handles the addition of a new courier
func AddCourier(c *gin.Context) {
	var courier models.Courier

	// Bind JSON request to the courier model
	if err := c.ShouldBindJSON(&courier); err != nil {
		log.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Set default availability if not provided
	if !courier.Available {
		courier.Available = true // Default to true
	}

	// Insert the courier into the database
	query := `INSERT INTO couriers (name, email, available) VALUES ($1, $2, $3) RETURNING id`
	var courierID int64
	err := db.DB.QueryRow(query, courier.Name, courier.Email, courier.Available).Scan(&courierID)
	if err != nil {
		log.Println("Error inserting courier into the database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add courier"})
		return
	}

	// Respond with the ID of the new courier
	c.JSON(http.StatusOK, gin.H{"message": "Courier added successfully", "courier_id": courierID})
}

// Example of another function: Get all couriers
func GetAllCouriers(c *gin.Context) {
	query := `SELECT * FROM couriers`
	var couriers []models.Courier

	err := db.DB.Select(&couriers, query)
	if err != nil {
		log.Println("Error fetching couriers from the database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch couriers"})
		return
	}

	c.JSON(http.StatusOK, couriers)
}

func GetAssignedOrdersByCourier(c *gin.Context) {
	// Retrieve user_id from the context
	// Retrieve user_id from the context (set by authentication middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		log.Println("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Query the database
	var orders []models.Order
	query := `SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status 
              FROM orders WHERE courier_id = $1`
	err := db.DB.Select(&orders, query, userID.(int64))
	if err != nil {
		log.Println("Error fetching assigned orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assigned orders"})
		return
	}

	log.Printf("Fetched %d assigned orders for courier ID: %d\n", len(orders), userID)

	// Return the orders
	c.JSON(http.StatusOK, orders)
}
