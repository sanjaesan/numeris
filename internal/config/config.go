package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// PostgresConfig -
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Dialect -
func (s PostgresConfig) Dialect() string {
	return "postgres"
}

// ConnectionInfo -
func (s PostgresConfig) ConnectionInfo() string {
	if s.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", s.Host, s.Port, s.User, s.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", s.Host, s.Port, s.User, s.Password, s.Name)
}

// DefaultPostgresConfig -
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		Name:     "numeris_invoice_dev_db",
	}
}

// Config represents the application configuration
type Config struct {
	LogFile    string         `json:"log_file"`
	CertFile   string         `json:"cert_file"`
	KeyFile    string         `json:"key_file"`
	Tls        bool           `json:"tls"`
	Port       int            `json:"port"`
	Database   PostgresConfig `json:"database"`
	CorsOrigin string         `json:"cors_origin"`
}

// DefaultConfig -
func DefaultConfig() Config {
	return Config{
		LogFile:    "",
		CertFile:   "",
		KeyFile:    "",
		Tls:        false,
		Port:       50056,
		Database:   DefaultPostgresConfig(),
		CorsOrigin: "",
	}
}

// LoadConfig -
func LoadConfig() Config {
	f, err := os.Open(".config")
	if err != nil {
		fmt.Println("Using the default config...")
		return DefaultConfig()
	}
	var c Config
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded .config")
	return c
}
