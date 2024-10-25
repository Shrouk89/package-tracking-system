package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Function that will handle user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Here you can add logic to process the registration form data (name, email, etc.)

	// Responding with a success message
	w.WriteHeader(http.StatusCreated) // 201 Created
	fmt.Fprintf(w, "User registration successful!")
}

func main() {
	// Initialize a new router
	router := mux.NewRouter()

	// Define the POST /register route
	router.HandleFunc("/register", RegisterUser).Methods("POST")

	// Start the server on port 8080
	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
