package main

import (
	"log"
)

type Runtime struct {
	Config   *Config
	Logger   *log.Logger
	Database *Database
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

	db, err := OpenDatabase(cfg)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Config:   cfg,
		Logger:   logger,
		Database: db,
	}, nil
}
