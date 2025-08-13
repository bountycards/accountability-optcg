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

	err = createGameResultsTable()
	if err != nil {
		return fmt.Errorf("failed to create game_results table: %w", err)
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

// createGameResultsTable creates the game_results table
func createGameResultsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS game_results (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		leader VARCHAR(100) NOT NULL,
		opponent VARCHAR(100) NOT NULL,
		category VARCHAR(50) DEFAULT 'Casual',
		went_first BOOLEAN NOT NULL,
		won BOOLEAN NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_game_results_user_id ON game_results(user_id);
	CREATE INDEX IF NOT EXISTS idx_game_results_leader ON game_results(leader);
	CREATE INDEX IF NOT EXISTS idx_game_results_category ON game_results(category);
	CREATE INDEX IF NOT EXISTS idx_game_results_created_at ON game_results(created_at);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create game_results table: %w", err)
	}

	// Add category column if it doesn't exist (for existing databases)
	err = addCategoryColumnIfNotExists()
	if err != nil {
		return fmt.Errorf("failed to add category column: %w", err)
	}

	log.Println("Game results table created successfully")
	return nil
}

// addCategoryColumnIfNotExists adds the category column to existing game_results tables
func addCategoryColumnIfNotExists() error {
	// Check if category column exists
	checkQuery := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = 'game_results' AND column_name = 'category'
	`

	var columnName string
	err := DB.QueryRow(checkQuery).Scan(&columnName)

	// If the column doesn't exist, add it
	if err == sql.ErrNoRows {
		alterQuery := `
			ALTER TABLE game_results 
			ADD COLUMN category VARCHAR(50) DEFAULT 'Casual'
		`

		_, err = DB.Exec(alterQuery)
		if err != nil {
			return fmt.Errorf("failed to add category column: %w", err)
		}

		// Add index for the new column
		indexQuery := `CREATE INDEX IF NOT EXISTS idx_game_results_category ON game_results(category)`
		_, err = DB.Exec(indexQuery)
		if err != nil {
			return fmt.Errorf("failed to create category index: %w", err)
		}

		log.Println("Added category column to existing game_results table")
	} else if err != nil {
		return fmt.Errorf("failed to check for category column: %w", err)
	}

	return nil
}
