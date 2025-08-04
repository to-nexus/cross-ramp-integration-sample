package config

import (
	"math/rand"
	"time"
)

// Config application configuration
type Config struct {
	Port string
	DB   DBConfig
}

// DBConfig database configuration
type DBConfig struct {
	Path string
}

// InitConfig initialize configuration
func InitConfig() *Config {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	return &Config{
		Port: ":8080",
		DB: DBConfig{
			Path: "./session_db",
		},
	}
}
