// internal/config/config.go
package config

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	logger "github.com/alfaa19/service-account-test/pkg/logrus"
	_ "github.com/lib/pq"
)

// Config holds application configuration
type Config struct {
	// Server settings (from command line args)
	Host string
	Port int

	// Database settings
	DB           *pgsql
	DBConnection *sql.DB

	// Logger settings
	LogLevel     logger.Level
	LogToConsole bool
	LogToFile    bool
	LogFilePath  string
}

// LoadConfig loads configuration from environment variables and command line arguments
func NewConfig() (*Config, error) {
	cfg := &Config{}

	// Parse command line arguments for server settings
	flag.StringVar(&cfg.Host, "host", "0.0.0.0", "Server host")
	flag.IntVar(&cfg.Port, "port", 8080, "Server port")

	flag.Parse()

	// Logger settings from environment variables
	cfg.LogLevel = logger.Level(strings.ToUpper(getEnv("LOG_LEVEL", "INFO")))
	cfg.LogToConsole = getEnv("LOG_CONSOLE", "true") == "true"
	cfg.LogToFile = getEnv("LOG_FILE", "true") == "true"
	cfg.LogFilePath = getEnv("LOG_PATH", "")

	// If log file path is empty and logging to file is enabled, set default path
	if cfg.LogFilePath == "" && cfg.LogToFile {
		// Get executable path
		execPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %v", err)
		}

		// Use logs directory in the same directory as the executable
		cfg.LogFilePath = filepath.Join(filepath.Dir(execPath), "logs", "application.log")
	}
	cfg.DB = newPostgres()
	return cfg, nil
}

// GetDBConnString returns the PostgreSQL connection string

// GetServerAddress returns the server address string
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
func (c *Config) OpenDatabase() error {
	if c.DBConnection == nil {
		db, err := c.DB.openPostgres()
		if err != nil {
			return err
		}
		c.DBConnection = db
	}
	return nil

}

// GetLoggerConfig returns the logger configuration
func (c *Config) GetLoggerConfig() logger.Config {
	return logger.Config{
		LogLevel:         c.LogLevel,
		LogToConsole:     c.LogToConsole,
		LogToFile:        c.LogToFile,
		LogFilePath:      c.LogFilePath,
		LogFileMaxSize:   10, // 10 MB
		LogFileMaxBackup: 5,  // Keep 5 backups
		LogFileMaxAge:    30, // 30 days
	}
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}
