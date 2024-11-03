package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

// InitDB initializes the database connection
func InitDB() {
	connStr := "user=postgres password=123 dbname=PackageTracking_db sslmode=disable"
	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database, try again: %v", err)
	}

	fmt.Println("Connected to the database, great!")
}
