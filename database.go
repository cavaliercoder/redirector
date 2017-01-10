package main

import (
	"fmt"
)

var (
	MappingNotFoundError = fmt.Errorf("Mapping not found")
)

type Database interface {
	Close() error
	AddMapping(m *Mapping) error
	GetMapping(key string) (*Mapping, error)
	GetMappings() ([]*Mapping, error)
	DeleteMapping(key string) error
	DeleteMappings() (int64, error)
	Stats() (DatabaseStats, error)
}

type DatabaseStats struct {
	TotalMappings int64 `json:"totalMappings"`
	DiskUsage     int64 `json:"diskUsage"`
}
