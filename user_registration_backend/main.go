package main

import (
	"fmt"
	"log"
	"net/http"
	"user_registration_backend/controllers" // ADDED
	"user_registration_backend/db" // ADDED
	"github.com/gorilla/mux"
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
	router := mux.NewRouter()

	// Define the POST /register route
	router.HandleFunc("/register", controllers.RegisterUser).Methods("POST") // Modified to use the function on controllers package

	// Define the POST /login route - ADDED: New route for login functionality
	router.HandleFunc("/login", controllers.LoginUser).Methods("POST")

	// Start the server on port 8080
	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
