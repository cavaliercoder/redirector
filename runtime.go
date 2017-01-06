package main

import (
	"log"
)

// Runtime contains globals for common runtime utilities.
type Runtime struct {
	Config       *Config
	Logger       *log.Logger
	AccessLogger *log.Logger
	Database     Database
}

func NewRuntime() (*Runtime, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger(cfg.LogFile)
	if err != nil {
		return nil, err
	}
	logger.Printf("Server started")

	accessLogger := logger
	if cfg.AccessLogFile != cfg.LogFile {
		accessLogger, err = NewLogger(cfg.AccessLogFile)
		if err != nil {
			return nil, err
		}
	}

	db, err := OpenBoltDatabase(cfg)
	if err != nil {
		return nil, err
	}
	logger.Printf("Using Bolt database: %v", cfg.DatabasePath)

	return &Runtime{
		Config:       cfg,
		Logger:       logger,
		AccessLogger: accessLogger,
		Database:     db,
	}, nil
}

func (rt *Runtime) Close() error {
	if rt.Database != nil {
		if err := rt.Database.Close(); err != nil {
			return err
		}
	}

	return nil
}
