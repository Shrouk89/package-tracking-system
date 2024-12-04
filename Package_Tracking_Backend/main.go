package main

import (
	"fmt"
	"net/http"

	"Package_Tracking_Backend/controllers" // ADDED
	"Package_Tracking_Backend/db"          // ADDED
	"Package_Tracking_Backend/middleware"

	"log"

	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// createSuperAdmin ensures a Super Admin account exists
func createSuperAdmin() {

	// Fetch the password from environment variables
	password := os.Getenv("SUPER_ADMIN_PASSWORD")
	if password == "" {
		log.Fatal("SUPER_ADMIN_PASSWORD environment variable is required")
	}

	// Check if a Super Admin already exists
	var count int
	query := `SELECT COUNT(*) FROM users WHERE role = $1`
	err := db.DB.Get(&count, query, "super_admin")
	if err != nil {
		log.Fatal("Failed to query database for Super Admin:", err)
	}

	if count > 0 {
		log.Println("Super Admin already exists. Skipping creation.")
		return
	}

	// Super Admin details (customize as needed)
	name := "Super Admin"
	email := "abdo21@gmail.com"

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	// Insert the Super Admin into the database
	query = `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4)`
	_, err = db.DB.Exec(query, name, email, string(hashedPassword), "super_admin")
	if err != nil {
		log.Fatal("Failed to insert Super Admin into database:", err)
	}

	log.Println("Super Admin created successfully with email:", email)
}

func main() {
	// Initialize database connection
	db.InitDB() // Ensures database connection is established

	createSuperAdmin()

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
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)

		// // Admin routes
		// public.POST("/admin/register", controllers.RegisterAdmin)
		// public.POST("/admin/login", controllers.LoginAdmin)
	}

	// Protected routes (authentication required)
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware to all routes in this group
	{
		protected.POST("/create-admin", controllers.RegisterAdmin)

		//protected.POST("/register", controllers.Register)
		protected.PUT("/approve-admin/:id", middleware.RoleMiddleware("super_admin"), controllers.ApproveAdmin)

		// Order routes
		protected.POST("/create-order", controllers.CreateOrder)
		protected.GET("/my-orders", controllers.GetOrdersByUser)
		protected.GET("/order-details/:id", controllers.GetOrderDetails)
		protected.PUT("/cancel-order/:id", controllers.CancelOrder)

		// Courier routes
		//protected.POST("/couriers", controllers.AddCourier)
		//protected.GET("/couriers", controllers.GetAllCouriers)
		protected.GET("/courier/assigned-orders", controllers.GetAssignedOrdersByCourier)
		protected.PUT("/courier/update-order-status/:id", controllers.UpdateOrderStatus)

		// Admin order management
		protected.GET("/admin/list-orders", controllers.GetAllOrders)
		protected.PUT("/admin/orders/update-status/:id", controllers.UpdateOrder)
		protected.DELETE("admin/delete-order/:id", controllers.DeleteOrder)
		protected.PUT("/admin/assign-order", controllers.AssignOrderToCourier)
		protected.GET("/admin/list-couriers", controllers.GetCurrentCouriers)

	}

	// Start the server on port 8080
	fmt.Println("Server is running at http://localhost:8080")
	router.Run(":8080")
}
