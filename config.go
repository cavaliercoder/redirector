package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// singleton config instance
var cfg *Config = &Config{
	Path:           "",
	Initialized:    false,
	DatabaseDriver: "bolt",
	DatabasePath:   "./redirector.db",
	ListenAddr:     ":8080",
	MgmtAddr:       "127.0.0.1:9321",
	LogFile:        "-", // stdout
	AccessLogFile:  "-", // stdout
	KeyBuilderName: "path",
}

// Config contains runtime configuration for the redirector service.
type Config struct {
	// Path of the loaded configuration file
	Path string `json:"-"`

	// Indicates if configuration has already been loaded
	Initialized bool `json:"-"`

	// Bypass panic handler for testing
	ExitOnError bool `json:"-"`

	// Database driver
	DatabaseDriver string `json:"database"`

	// File path of the BoltDB database
	DatabasePath string `json:"databasePath"`

	// Interface and TCP port to bind to
	ListenAddr string `json:"listenAddr"`

	// Interface and TCP to bind the management service to
	MgmtAddr string `json:"mgmtAddr"`

	// Logfile to write to
	LogFile string `json:"logFile"`

	// LogFile to write HTTP transactions to
	AccessLogFile string `json:"accessLogFile"`

	// The name of the key build to use when mapping URLs
	KeyBuilderName string `json:"keyBuilder"`

	// An instance of a KeyBuilder to use when mapping URLs
	KeyBuilder KeyBuilder `json:"-"`

	// Use this key if the requested key does not exist (instead of returning
	// 404)
	DefaultKey string `json:"defaultKey"`
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
	if err := cfg.initialize(); err != nil {
		return nil, err
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

		if !cfg.Initialized {
			if err := cfg.initialize(); err != nil {
				return nil, err
			}
		}
	}

	return cfg, nil
}

// initialize instanciates the current runtime configuration.
func (c *Config) initialize() error {
	// expand keybuilder
	switch c.KeyBuilderName {
	case "", "path":
		c.KeyBuilder = RequestURIPathKeyBuilder()

	case "uri":
		c.KeyBuilder = RequestURIKeyBuilder()

	default:
		if strings.HasPrefix(c.KeyBuilderName, "param:") {
			c.KeyBuilder = RequestParamKeyBuilder(c.KeyBuilderName[6:])
		} else {
			return fmt.Errorf("Unknown key builder: %v", c.KeyBuilderName)
		}
	}

	if err := InitTemplates(); err != nil {
		return err
	}

	c.Initialized = true
	return nil
}
