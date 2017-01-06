package main

import (
	"testing"
)

var testMappings = []Mapping{
	{"default", "/okay", false},
	{"/permanent", "/okay", true},
	{"/temporary", "/okay", false},
}

func testDB(t *testing.T, db Database) {
	// add mappings
	for _, m := range testMappings {
		if err := db.AddMapping(&m); err != nil {
			panic(err)
		}
	}

	// test mappings
	for _, m := range testMappings {
		if v, err := db.GetMapping(m.Key); err != nil {
			panic(err)
		} else {
			if v.Destination != m.Destination {
				t.Errorf("Bad mapping destination: '%v', expected '%v'", v.Destination, m.Destination)
			}

			if v.Permanent != m.Permanent {
				t.Errorf("Bad mapping permanence: '%v', expected '%v'", v.Permanent, m.Permanent)
			}
		}
	}

	// test all
	if v, err := db.GetMappings(); err != nil {
		panic(err)
	} else {
		if len(v) != len(testMappings) {
			t.Errorf("Bad mapping count %v, expected %v", len(v), len(testMappings))
		}
	}

	// delete mappings
	for _, m := range testMappings {
		if err := db.DeleteMapping(m.Key); err != nil {
			panic(err)
		}

		if _, err := db.GetMapping(m.Key); err != MappingNotFoundError {
			t.Errorf("Mapping was deleted but still exists in database")
		}
	}
}
