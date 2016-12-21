package main

import (
	"log"
	"os"
)

func NewLogger(cfg *Config) (*log.Logger, error) {
	out := os.Stdout
	if cfg.LogFile != "-" {
		f, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
		if err != nil {
			return nil, err
		}

		out = f
	}

	return log.New(out, "", log.Ldate|log.Ltime|log.Lmicroseconds), nil
}
