package main

import (
	"fmt"
	"net/http"

	"Package_Tracking_Backend/controllers" // ADDED
	"Package_Tracking_Backend/db"          // ADDED
	"Package_Tracking_Backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	db.InitDB() // Ensures database connection is established

	// Initialize a new router
	router := gin.Default()

	// Enable CORS to allow the Angular frontend to communicate with the Go backend
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

	// Public routes (accessible without authentication)
	public := router.Group("/")
	{
		// User routes
		public.POST("/register", controllers.RegisterUser)
		public.POST("/login", controllers.LoginUser)

		// Admin routes
		public.POST("/admin/register", controllers.RegisterAdmin)
		public.POST("/admin/login", controllers.LoginAdmin)
	}

	// Protected routes (authentication required)
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware to all routes in this group
	{
		// Order routes
		protected.POST("/create-order", controllers.CreateOrder)
		protected.GET("/my-orders", controllers.GetOrdersByUser)
		protected.GET("/order-details/:id", controllers.GetOrderDetails)
		protected.PUT("/cancel-order/:id", controllers.CancelOrder)

		// Courier routes
		protected.POST("/couriers", controllers.AddCourier)
		protected.GET("/couriers", controllers.GetAllCouriers)
		protected.GET("/assigned-orders", controllers.GetAssignedOrdersByCourier)
		protected.PUT("/orders/update-status/:id", controllers.UpdateOrderStatus)

		// Admin order management
		protected.GET("/orders", controllers.GetAllOrders)
		protected.PUT("/orders/:id", controllers.UpdateOrder)
		protected.DELETE("/orders/:id", controllers.DeleteOrder)
	}

	// Start the server on port 8080
	fmt.Println("Server is running at http://localhost:8080")
	router.Run(":8080")
}
