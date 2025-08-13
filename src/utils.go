package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func getEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		value = os.Getenv(key)
	}
	return value
}
