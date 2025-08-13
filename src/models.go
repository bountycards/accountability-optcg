package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// User represents a Discord user in our system
type User struct {
	ID            int       `json:"id"`
	DiscordID     string    `json:"discord_id"`
	Username      string    `json:"username"`
	Discriminator string    `json:"discriminator"`
	Timezone      string    `json:"timezone"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateUser inserts a new user into the database
func CreateUser(discordID, username, discriminator string) (*User, error) {
	query := `
		INSERT INTO users (discord_id, username, discriminator)
		VALUES ($1, $2, $3)
		RETURNING id, discord_id, username, discriminator, timezone, created_at, updated_at
	`

	user := &User{}
	err := DB.QueryRow(query, discordID, username, discriminator).Scan(
		&user.ID,
		&user.DiscordID,
		&user.Username,
		&user.Discriminator,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByDiscordID retrieves a user by their Discord ID
func GetUserByDiscordID(discordID string) (*User, error) {
	query := `
		SELECT id, discord_id, username, discriminator, timezone, created_at, updated_at
		FROM users
		WHERE discord_id = $1
	`

	user := &User{}
	err := DB.QueryRow(query, discordID).Scan(
		&user.ID,
		&user.DiscordID,
		&user.Username,
		&user.Discriminator,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUserTimezone updates a user's timezone
func UpdateUserTimezone(discordID, timezone string) error {
	query := `
		UPDATE users 
		SET timezone = $1, updated_at = NOW()
		WHERE discord_id = $2
	`

	result, err := DB.Exec(query, timezone, discordID)
	if err != nil {
		return fmt.Errorf("failed to update user timezone: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetOrCreateUser gets an existing user or creates a new one
func GetOrCreateUser(discordID, username, discriminator string) (*User, error) {
	// Try to get existing user first
	user, err := GetUserByDiscordID(discordID)
	if err == nil {
		return user, nil
	}

	// If user doesn't exist, create a new one
	if err.Error() == "user not found" {
		return CreateUser(discordID, username, discriminator)
	}

	// Some other error occurred
	return nil, fmt.Errorf("failed to get or create user: %w", err)
}

// GameResult represents a game result record
type GameResult struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Leader    string    `json:"leader"`
	Opponent  string    `json:"opponent"`
	Category  string    `json:"category"`
	WentFirst bool      `json:"went_first"`
	Won       bool      `json:"won"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateGameResult inserts a new game result into the database
func CreateGameResult(userID int, leader, opponent, category string, wentFirst, won bool) (*GameResult, error) {
	query := `
		INSERT INTO game_results (user_id, leader, opponent, category, went_first, won)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, leader, opponent, category, went_first, won, created_at
	`

	gameResult := &GameResult{}
	err := DB.QueryRow(query, userID, leader, opponent, category, wentFirst, won).Scan(
		&gameResult.ID,
		&gameResult.UserID,
		&gameResult.Leader,
		&gameResult.Opponent,
		&gameResult.Category,
		&gameResult.WentFirst,
		&gameResult.Won,
		&gameResult.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create game result: %w", err)
	}

	return gameResult, nil
}

// ValidateCategory checks if the provided category is valid
func ValidateCategory(category string) bool {
	validCategories := []string{
		"Casual",
		"Ranked",
		"Locals",
		"Regional",
		"National",
		"Tournament",
		"Practice",
		"Online",
	}

	for _, valid := range validCategories {
		if strings.EqualFold(category, valid) {
			return true
		}
	}

	return false
}

// NormalizeCategory normalizes the category to proper capitalization
func NormalizeCategory(category string) string {
	validCategories := map[string]string{
		"casual":     "Casual",
		"ranked":     "Ranked",
		"locals":     "Locals",
		"regional":   "Regional",
		"national":   "National",
		"tournament": "Tournament",
		"practice":   "Practice",
		"online":     "Online",
	}

	normalized, exists := validCategories[strings.ToLower(category)]
	if exists {
		return normalized
	}

	return "Casual" // Default fallback
}
