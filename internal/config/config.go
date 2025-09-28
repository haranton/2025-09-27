package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	StoragePath string `env:"STORAGEPATH" envDefault:"storage"`
	MaxWorkers  int    `env:"MAXWORKERS" envDefault:"4"`
	HTTPPort    int    `env:"PORT" envDefault:"8080"`
}

var AppConfig *Config

func InitConfig() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	AppConfig = cfg
	return nil
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%d", c.HTTPPort)
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
