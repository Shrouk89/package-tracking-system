package main

import (
	"fmt"
	"log"
	"net/http"
	"user_registration_backend/controllers" // ADDED
	"user_registration_backend/db" // ADDED
	"github.com/gin-gonic/gin"
        "user_registration_backend/controllers" //use gin instead of mux
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
	db.InitDB()  // Ensures database connection is established
	
	// Initialize a new router
    router := gin.Default()

    // Define the POST /register route
    router.POST("/register", controllers.RegisterUser)

    // Define the POST /login route
    router.POST("/login", controllers.LoginUser)

    // Start the server on port 8080
    fmt.Println("Server is running at http://localhost:8080")
    router.Run(":8080")
}
