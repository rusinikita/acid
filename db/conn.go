package db

import (
	"log"
	"os"
	"time"

	"database/sql"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	// Load .env file from the root directory
	if err := godotenv.Load(".env"); err != nil {
		//log.Println(".env not found")
	}
	if err := godotenv.Load("../.env"); err != nil {
		//log.Println("../.env not found")
	}
	if err := godotenv.Load("../../.env"); err != nil {
		//log.Println("../../.env not found")
	}

	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_CONNECT"))
	if err != nil {
		log.Fatalf("Failed to Connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to Ping to database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}
