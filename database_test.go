package main

import (
	"testing"
)

var testMappings = []Mapping{
	{"default", "/okay", false, "Should only apply to missing keys", false},
	{"/permanent", "/okay", true, "Should return HTTP 308", false},
	{"/temporary", "/okay", false, "Should return HTTP 307", false},
	{"/template", "/?key={{ .Key }}", false, "Should expand template", true},
}

func testDB(t *testing.T, db Database) {
	// clear existing
	if _, err := db.DeleteMappings(); err != nil {
		panic(err)
	}

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

			if v.Comment != m.Comment {
				t.Errorf("Bad mapping comment: '%v', expected '%v'", v.Comment, m.Comment)
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
