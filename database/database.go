package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// Get the environment credentials and opens a connection with the database
func Connect() *sql.DB {
	password := os.Getenv("PG_PASS")
	user := os.Getenv("PG_USER")
	dbName := os.Getenv("PG_DB_STORE")
	host := os.Getenv("PG_HOST")
	connection := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbName, password, host)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Println("Unable to connect:" + err.Error())
	}
	return db
}
