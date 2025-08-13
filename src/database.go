package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection with retry logic
func InitDB() error {
	dbHost := getEnv("POSTGRES_HOST")
	dbPort := getEnv("POSTGRES_PORT")
	dbUser := getEnv("POSTGRES_USER")
	dbPassword := getEnv("POSTGRES_PASSWORD")
	dbName := getEnv("POSTGRES_DB")

	// Default values if not set
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Get retry configuration from environment variables with defaults
	maxRetries := 5
	retryDelaySeconds := 10

	if retriesStr := getEnv("DB_MAX_RETRIES"); retriesStr != "" {
		if retries, err := strconv.Atoi(retriesStr); err == nil && retries > 0 {
			maxRetries = retries
		}
	}

	if delayStr := getEnv("DB_RETRY_DELAY_SECONDS"); delayStr != "" {
		if delay, err := strconv.Atoi(delayStr); err == nil && delay > 0 {
			retryDelaySeconds = delay
		}
	}

	retryDelay := time.Duration(retryDelaySeconds) * time.Second

	var err error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Attempting to connect to database (attempt %d/%d)...", attempt, maxRetries)

		DB, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Printf("Failed to open database connection: %v", err)
			if attempt == maxRetries {
				return fmt.Errorf("failed to open database after %d attempts: %w", maxRetries, err)
			}
			time.Sleep(retryDelay)
			continue
		}

		err = DB.Ping()
		if err != nil {
			log.Printf("Failed to ping database: %v", err)
			DB.Close() // Close the connection before retrying
			if attempt == maxRetries {
				return fmt.Errorf("failed to ping database after %d attempts: %w", maxRetries, err)
			}
			log.Printf("Waiting %v before next attempt...", retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		// Success!
		log.Println("Successfully connected to PostgreSQL database")
		return nil
	}

	return fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// CreateTables creates all necessary tables
func CreateTables() error {
	err := createUsersTable()
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Successfully created all database tables")
	return nil
}

// createUsersTable creates the users table
func createUsersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		discord_id VARCHAR(20) UNIQUE NOT NULL,
		username VARCHAR(32) NOT NULL,
		discriminator VARCHAR(4),
		timezone VARCHAR(50) DEFAULT 'UTC',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_users_discord_id ON users(discord_id);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Users table created successfully")
	return nil
}
