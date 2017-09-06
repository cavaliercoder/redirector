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
	ViewBag:        NewViewBag(),
}

// Config contains runtime configuration for the redirector service.
type Config struct {
	Path              string     `json:"-"` // Loaded configuration file
	Initialized       bool       `json:"-"`
	ExitOnError       bool       `json:"-"` // Bypass panic handler
	DatabaseDriver    string     `json:"database"`
	DatabasePath      string     `json:"databasePath"`
	ListenAddr        string     `json:"listenAddr"`
	MgmtAddr          string     `json:"mgmtAddr"`
	LogFile           string     `json:"logFile"`
	AccessLogFile     string     `json:"accessLogFile"`
	KeyBuilderName    string     `json:"keyBuilder"` // The name of the KeyBuilder
	KeyBuilder        KeyBuilder `json:"-"`          // An instance of a KeyBuilder
	DefaultKey        string     `json:"defaultKey"` // fallback for all 404s
	DestinationPrefix string     `json:"destinationPrefix"`
	ViewBag           ViewBag    `json:"viewBag"`
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
