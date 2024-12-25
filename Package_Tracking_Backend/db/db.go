package db

import (
	"fmt"
	"log"
	"os"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

// The name of the db is changed to include the whole data not only the users
// InitDB initializes the database connection
func InitDB() {
	//connStr := "user=postgres password=123 dbname=PackageTracking_db sslmode=disable"
	//var err error
	//DB, err = sqlx.Connect("postgres", connStr)
	//if err != nil {
	//	log.Fatalf("Error connecting to the database, try again: %v", err)
	//}

	//fmt.Println("Connected to the database, great!")

	dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    // Build the connection string
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName,
    )

    var err error
    DB, err = sqlx.Connect("postgres", connStr)
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }

    fmt.Println("Connected to the database, great!")

}
