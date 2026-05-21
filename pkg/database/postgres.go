package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Connect loads .env, builds the connection string, and returns an open DB handle.
// The caller is responsible for closing the DB when the application shuts down.
func Connect() (*sqlx.DB, error) {
	if err := godotenv.Load(); err != nil { 
		return nil, fmt.Errorf("could not load .env file: %w", err)
	}
	//read the .env file and build the connection string dsn = "Data Source Name"
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return db, nil
}
