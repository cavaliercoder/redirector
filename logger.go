package main

import (
	"log"
	"os"
)

func NewLogger(path string) (*log.Logger, error) {
	out := os.Stdout
	if path != "-" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
		if err != nil {
			return nil, err
		}

		out = f
	}

	return log.New(out, "", log.Ldate|log.Ltime|log.Lmicroseconds), nil
}
