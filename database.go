package main

import (
	"fmt"
)

var (
	MappingNotFoundError = fmt.Errorf("Mapping not found")
)

type Database interface {
	Close() error
	Stats() interface{}
	AddMapping(m *Mapping) error
	GetMapping(key string) (*Mapping, error)
	GetMappings() ([]*Mapping, error)
	DeleteMapping(key string) error
}
