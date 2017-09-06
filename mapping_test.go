package main

import (
	"testing"
)

func TestMappingValidation(t *testing.T) {
	m := &Mapping{}
	if err := m.Validate(); err != KeyMissingError {
		t.Errorf("Expected mapping validation to fail with %T, got '%v'", KeyMissingError, err)
	}

	m.Key = "/test"
	if err := m.Validate(); err != DestinationMissingError {
		t.Errorf("Expected mapping validation to fail with %T, got '%v'", DestinationMissingError, err)
	}

	m.Destination = "/?key={{ .Key }}"
	if err := m.Validate(); err != nil {
		t.Fatalf("Expected validation to pass, got '%v'", err)
	}

	if !m.IsTemplate {
		t.Fatalf("Expect Mapping.IsTemplate to be true, got false")
	}
}

func TestDestinationTemplate(t *testing.T) {
	m := &Mapping{
		Key:         "/template",
		Destination: "/?key={{ .Key }}",
	}
	expect := "/?key=/template"

	if err := m.Validate(); err != nil {
		panic(err)
	}

	dest, err := m.ComputeDestination(ViewBag{"Key": m.Key})
	if err != nil {
		panic(err)
	}

	if dest != expect {
		t.Fatalf("Expected destination '%v', got '%v'", expect, dest)
	}
}
