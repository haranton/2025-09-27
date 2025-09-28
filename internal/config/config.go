package config

import (
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	StoragePath string `env:"STORAGEPATH" envDefault:"storage"`
	MaxWorkers  int    `env:"MAXWORKERS" envDefault:"4"`
	HTTPPort    int    `env:"PORT" envDefault:"8080"`
}

func getConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Note: No .env file found")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return &cfg, nil

}
