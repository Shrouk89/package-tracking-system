package main

import (
	"fmt"
	"net/http"
	//"log"
	//"net/http"
	"Package_Tracking_Backend/controllers" // ADDED
	"Package_Tracking_Backend/db"          // ADDED

	"github.com/gin-gonic/gin"
	//"Package_Tracking_Backend/controllers" //use gin instead of mux
)

// Commented: I think we can just use the already implemented function from the controllers package.

// Function that will handle user registration
//func RegisterUser(w http.ResponseWriter, r *http.Request) {
// Here you can add logic to process the registration form data (name, email, etc.)

// Responding with a success message
//	w.WriteHeader(http.StatusCreated) // 201 Created
//	fmt.Fprintf(w, "User registration successful!")
//}

func main() {
	// Initialize database connection - ADDED
	db.InitDB() // Ensures database connection is established

	// Initialize a new router
	router := gin.Default()

	// Enable CORS to allow our Angular frontend to communicate with this Go backend
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // Allow all origins
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Define the POST /register endpoint
	router.POST("/register", controllers.RegisterUser)

	// Define the POST /login endpoint
	router.POST("/login", controllers.LoginUser)

	// Define the POST /create-order endpoint
	router.POST("/create-order", controllers.CreateOrder)

	// Define the GET /my-orders endpoint
	router.GET("/my-orders", controllers.GetOrdersByUser)

	// Define the GET /order-details by ID endpoint
	router.GET("/order-details/:id", controllers.GetOrderDetails)

	// Start the server on port 8080
	fmt.Println("Server is running at http://localhost:8080")
	router.Run(":8080")
}
