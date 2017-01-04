package main

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	if cfg, err := GetConfig(); err != nil {
		panic(err)
	} else {
		if cfg.Path != "" {
			t.Fatalf("Expected default config, got %v", cfg.Path)
		}
	}
}
