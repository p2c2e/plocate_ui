package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type IndexConfig struct {
	Name         string   `yaml:"name"`
	DatabasePath string   `yaml:"database_path"`
	IndexPaths   []string `yaml:"index_paths"`
	Enabled      bool     `yaml:"enabled"`
}

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Plocate struct {
		// Legacy fields for backward compatibility
		DatabasePath string   `yaml:"database_path,omitempty"`
		IndexPaths   []string `yaml:"index_paths,omitempty"`

		// New multi-index configuration
		Indices     []IndexConfig `yaml:"indices,omitempty"`
		UpdatedbBin string        `yaml:"updatedb_bin"`
		PlocateBin  string        `yaml:"plocate_bin"`
	} `yaml:"plocate"`

	Scheduler struct {
		Enabled  bool   `yaml:"enabled"`
		Interval string `yaml:"interval"` // cron format: "0 */6 * * *" = every 6 hours
	} `yaml:"scheduler"`
}

var (
	AppConfig  *Config
	configPath string
	mu         sync.Mutex
)

func Load(path string) error {
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
		if path == "" {
			path = "/app/config/config.yml"
		}
	}
	configPath = path

	var cfg Config

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		// No config file — start with empty defaults
		cfg = defaultConfig()
	} else if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	} else {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}
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
	if cfg.Plocate.UpdatedbBin == "" {
		cfg.Plocate.UpdatedbBin = "updatedb"
	}
	if cfg.Plocate.PlocateBin == "" {
		cfg.Plocate.PlocateBin = "plocate"
	}
	if cfg.Scheduler.Interval == "" {
		cfg.Scheduler.Interval = "0 */6 * * *" // Every 6 hours by default
	}

	// Handle backward compatibility: convert old format to new format
	if len(cfg.Plocate.Indices) == 0 && cfg.Plocate.DatabasePath != "" {
		defaultIndex := IndexConfig{
			Name:         "default",
			DatabasePath: cfg.Plocate.DatabasePath,
			IndexPaths:   cfg.Plocate.IndexPaths,
			Enabled:      true,
		}
		if len(defaultIndex.IndexPaths) == 0 {
			defaultIndex.IndexPaths = []string{"/"}
		}
		cfg.Plocate.Indices = []IndexConfig{defaultIndex}
	}

	// Ensure all index database directories exist
	for _, index := range cfg.Plocate.Indices {
		dbDir := filepath.Dir(index.DatabasePath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory for index %s: %w", index.Name, err)
		}
	}

	AppConfig = &cfg

	// Persist the config so it exists on disk for next startup
	if os.IsNotExist(err) {
		if saveErr := Save(); saveErr != nil {
			return fmt.Errorf("failed to save initial config: %w", saveErr)
		}
	}

	return nil
}

func defaultConfig() Config {
	var cfg Config
	cfg.Server.Port = "8080"
	cfg.Plocate.UpdatedbBin = "updatedb"
	cfg.Plocate.PlocateBin = "plocate"
	cfg.Scheduler.Enabled = true
	cfg.Scheduler.Interval = "0 */6 * * *"
	return cfg
}

// Save persists the current AppConfig to the config file.
func Save() error {
	mu.Lock()
	defer mu.Unlock()

	data, err := yaml.Marshal(AppConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddIndex adds a new index to the config and persists it.
func AddIndex(name string, paths []string) (*IndexConfig, error) {
	mu.Lock()
	defer mu.Unlock()

	// Check for duplicate name
	for _, idx := range AppConfig.Plocate.Indices {
		if idx.Name == name {
			return nil, fmt.Errorf("index '%s' already exists", name)
		}
	}

	dbPath := filepath.Join("/app/data", name+".db")

	idx := IndexConfig{
		Name:         name,
		DatabasePath: dbPath,
		IndexPaths:   paths,
		Enabled:      true,
	}

	// Ensure database directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	AppConfig.Plocate.Indices = append(AppConfig.Plocate.Indices, idx)

	// Save without holding mu (Save has its own lock) — release and re-save
	// Actually, Save also locks mu, so we need to save inline here.
	data, err := yaml.Marshal(AppConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	return &idx, nil
}

// RemoveIndex removes an index from the config and persists it.
func RemoveIndex(name string) error {
	mu.Lock()
	defer mu.Unlock()

	found := false
	indices := make([]IndexConfig, 0, len(AppConfig.Plocate.Indices))
	for _, idx := range AppConfig.Plocate.Indices {
		if idx.Name == name {
			found = true
			continue
		}
		indices = append(indices, idx)
	}

	if !found {
		return fmt.Errorf("index '%s' not found", name)
	}

	AppConfig.Plocate.Indices = indices

	// Save inline (same mutex)
	data, err := yaml.Marshal(AppConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
