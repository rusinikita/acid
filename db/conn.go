package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// LoadEnv loads .env from baseDir and its parent directory.
func LoadEnv(baseDir string) {
	_ = godotenv.Load(filepath.Join(baseDir, ".env"))
	_ = godotenv.Load(filepath.Join(filepath.Dir(baseDir), ".env"))
}

// TryConnect attempts to open and ping the database, returning an error on failure.
func TryConnect() error {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	database, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_CONNECT"))
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	return nil
}

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
