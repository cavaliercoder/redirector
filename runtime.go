package main

import (
	"fmt"
	"log"
)

// Runtime contains globals for common runtime utilities.
type Runtime struct {
	Config       *Config
	Logger       *log.Logger
	AccessLogger *log.Logger
	Database     Database
}

var UnsupportedDatabaseDriverError = fmt.Errorf("Unsupported database driver")

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

	var db Database
	switch cfg.DatabaseDriver {
	case "bolt":
		db, err = OpenBoltDatabase(cfg)

	case "redis":
		db, err = OpenRedisDatabase(cfg)
	default:
		return nil, UnsupportedDatabaseDriverError
	}

	if err != nil {
		return nil, err
	}

	dbstats, err := db.Stats()
	if err != nil {
		return nil, err
	}

	logger.Printf("Connected to database %v://%v", cfg.DatabaseDriver, cfg.DatabasePath)
	logger.Printf("  Total mappings: %v\n", dbstats.TotalMappings)
	logger.Printf("  Disk usage: %v bytes\n", dbstats.DiskUsage)

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
