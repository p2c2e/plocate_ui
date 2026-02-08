package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Plocate struct {
		DatabasePath string   `yaml:"database_path"`
		IndexPaths   []string `yaml:"index_paths"`
		UpdatedbBin  string   `yaml:"updatedb_bin"`
		PlocateBin   string   `yaml:"plocate_bin"`
	} `yaml:"plocate"`

	Scheduler struct {
		Enabled  bool   `yaml:"enabled"`
		Interval string `yaml:"interval"` // cron format: "0 */6 * * *" = every 6 hours
	} `yaml:"scheduler"`
}

var AppConfig *Config

func Load(configPath string) error {
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			configPath = "/app/config.yml"
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}
	if dbPath := os.Getenv("PLOCATE_DB_PATH"); dbPath != "" {
		cfg.Plocate.DatabasePath = dbPath
	}
	if interval := os.Getenv("INDEX_INTERVAL"); interval != "" {
		cfg.Scheduler.Interval = interval
	}

	// Set defaults
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Plocate.DatabasePath == "" {
		cfg.Plocate.DatabasePath = "/var/lib/plocate/plocate.db"
	}
	if cfg.Plocate.UpdatedbBin == "" {
		cfg.Plocate.UpdatedbBin = "updatedb"
	}
	if cfg.Plocate.PlocateBin == "" {
		cfg.Plocate.PlocateBin = "plocate"
	}
	if cfg.Scheduler.Interval == "" {
		cfg.Scheduler.Interval = "0 */6 * * *" // Every 6 hours by default
	}

	// Ensure database directory exists
	dbDir := filepath.Dir(cfg.Plocate.DatabasePath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	AppConfig = &cfg
	return nil
}
