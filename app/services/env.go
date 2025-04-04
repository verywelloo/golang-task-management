package service

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func EnvMongoURI() string {
	err := godotenv.Load("env.sh")
	if err != nil {
		log.Fatalf("Error loading env file  %v", err)
		//LogDetail(6, "Error loading env file")
	}

	return GetEnv("MONGOURI", "")
}
