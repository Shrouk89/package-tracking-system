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
	// Retrieve user_id from the context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		log.Println("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		log.Println("Failed to assert user ID type:", userIDInterface)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	log.Printf("Fetching assigned orders for courier ID: %d\n", userID)

	// Query the database
	var orders []models.Order
	query := `SELECT id, user_id, pickup_location, dropoff_location, package_details, delivery_time, status 
              FROM orders WHERE courier_id = $1`
	err := db.DB.Select(&orders, query, userID)
	if err != nil {
		log.Println("Error fetching assigned orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assigned orders"})
		return
	}

	log.Printf("Fetched %d assigned orders for courier ID: %d\n", len(orders), userID)

	// Return the orders
	c.JSON(http.StatusOK, orders)
}

// UpdateOrderStatus handles updating the status of an order

func UpdateOrderStatus(c *gin.Context) {
	courierID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Courier not authenticated"})
		return
	}

	orderID := c.Param("id")
	var requestBody struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	validStatuses := map[string]bool{
		"Picked Up":  true,
		"In Transit": true,
		"Delivered":  true,
	}
	if !validStatuses[requestBody.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status value"})
		return
	}

	// Ensure the order belongs to the authenticated courier
	queryCheck := `SELECT id FROM orders WHERE id = $1 AND courier_id = $2`
	var id int64
	err := db.DB.Get(&id, queryCheck, orderID, courierID.(int64))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Order not assigned to this courier"})
		return
	}

	// Update order status
	queryUpdate := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err = db.DB.Exec(queryUpdate, requestBody.Status, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// rowsAffected, err := result.RowsAffected()
	// if err != nil {
	// 	log.Println("Error fetching rows affected:", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
	// 	return
	// }

	// if rowsAffected == 0 {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}
