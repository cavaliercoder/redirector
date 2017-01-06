package main

import (
	"testing"
)

func TestTemplates(t *testing.T) {
	if err := InitTemplates(); err != nil {
		panic(err)
	}

	for _, code := range statusCodes {
		body, err := BodyForStatus(code)
		if err != nil {
			panic(err)
		}

		if body == "" {
			t.Fatalf("Empty body returned for status %v", code)
		}
	}
}
