package config

import (
	"log"
	"os"
)

type Config struct {
	Port   string
	APIKey string
}

func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("Warning: API_KEY not set. Using default value for development.")
		apiKey = "AIzaSyCZCRGPPeUlacQ2bMap9EsffHrofCuyrko"
	}

	return &Config{
		Port:   port,
		APIKey: apiKey,
	}
}
