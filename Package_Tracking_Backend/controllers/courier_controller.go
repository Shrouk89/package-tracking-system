package controllers

import (
	"log"
	"net/http"

	"Package_Tracking_Backend/db"
	"Package_Tracking_Backend/models"

	"github.com/gin-gonic/gin"
)

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
	courierID := c.Query("courier_id") // Get courier ID from query parameters

	if courierID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Courier ID is required"})
		return
	}

	var orders []models.Order
	query := `SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status 
              FROM orders WHERE courier_id = $1` // Query to get orders by courier ID
	err := db.DB.Select(&orders, query, courierID)

	if err != nil {
		log.Println("Error fetching assigned orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assigned orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus handles updating the status of an order
func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id") // Get the order ID from the URL parameters
	var requestBody struct {
		Status string `json:"status" binding:"required"` // Require a new status in the request body
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Prepare the query to update the order status
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := db.DB.Exec(query, requestBody.Status, orderID)

	if err != nil {
		log.Println("Error updating order status:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
