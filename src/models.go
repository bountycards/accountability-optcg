package main

import (
	"database/sql"
	"fmt"
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
