package main

import (
	"log"
)

// Runtime contains globals for common runtime utilities.
type Runtime struct {
	Config   *Config
	Logger   *log.Logger
	Database Database
}

func NewRuntime() (*Runtime, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		return nil, err
	}

	db, err := OpenBoltDatabase(cfg)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Config:   cfg,
		Logger:   logger,
		Database: db,
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
