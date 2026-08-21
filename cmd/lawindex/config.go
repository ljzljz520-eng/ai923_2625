package main

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	DatabasePath string
	Address      string
}

func DefaultConfig() Config { return Config{DatabasePath: "lawindex.db", Address: ":8080"} }
func ConfigFromEnv() Config {
	config := DefaultConfig()
	if value := strings.TrimSpace(os.Getenv("LAWINDEX_DB")); value != "" {
		config.DatabasePath = value
	}
	if value := strings.TrimSpace(os.Getenv("LAWINDEX_ADDR")); value != "" {
		config.Address = value
	}
	return config
}
func (c Config) Validate() error {
	if strings.TrimSpace(c.DatabasePath) == "" {
		return errors.New("database path is required")
	}
	if strings.TrimSpace(c.Address) == "" {
		return errors.New("address is required")
	}
	return nil
}
func (c Config) IsUnixSocket() bool { return strings.HasPrefix(c.Address, "/") }
func (c Config) Display() string    { return c.DatabasePath + " @ " + c.Address }
