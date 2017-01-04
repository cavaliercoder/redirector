package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// singleton config instance
var cfg *Config = &Config{
	Path:           "",
	Initialized:    false,
	DatabasePath:   "./redirector.db",
	ListenAddr:     ":8080",
	MgmtAddr:       "127.0.0.1:9321",
	LogFile:        "-", // stdout
	KeyBuilderName: "path",
}

// Config contains runtime configuration for the redirector service.
type Config struct {
	// Path of the loaded configuration file
	Path string `json:"-"`

	// Indicates if configuration has already been loaded
	Initialized bool `json:"-"`

	// File path of the BoltDB database
	DatabasePath string `json:"databasePath"`

	// Interface and TCP port to bind to
	ListenAddr string `json:"listenAddr"`

	// Interface and TCP to bind the management service to
	MgmtAddr string `json:"mgmtAddr"`

	// Logfile to write to
	LogFile string `json:"logFile"`

	// The name of the key build to use when mapping URLs
	KeyBuilderName string `json:"keyBuilder"`

	// An instance of a KeyBuilder to use when mapping URLs
	KeyBuilder KeyBuilder `json:"-"`
}

// LoadConfig reads configuration from the given file path
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	if err := dec.Decode(cfg); err != nil {
		return nil, err
	}

	cfg.Path = path
	cfg.Initialized = true

	// expand keybuilder
	// TODO: add query param key builder
	switch cfg.KeyBuilderName {
	case "", "path":
		cfg.KeyBuilder = RequestURIPathKeyBuilder()

	case "uri":
		cfg.KeyBuilder = RequestURIKeyBuilder()

	default:
		return nil, fmt.Errorf("Unknown key builder: %v", cfg.KeyBuilderName)
	}

	return cfg, nil
}

// GetConfig returns a pointer to a singleton configuration struct.
func GetConfig() (*Config, error) {
	if !cfg.Initialized {
		cfgPaths := []string{
			"./redirector.json",
			"/etc/redirector/redirector.json",
		}

		for _, path := range cfgPaths {
			_, err := LoadConfig(path)
			if err == nil {
				break
			}

			if !os.IsNotExist(err) {
				return nil, err
			}
		}
	}

	return cfg, nil
}
