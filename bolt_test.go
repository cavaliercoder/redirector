package main

import (
	"io/ioutil"
	"os"
	"testing"
)

var boltdbMappings = []Mapping{
	{"/test1", "/yes!", true},
	{"/test2", "/nope", false},
}

func TestBoltDBMappings(t *testing.T) {
	// create temp file
	dbpath := func() string {
		if dbf, err := ioutil.TempFile("", "boltdb_test_"); err != nil {
			panic(err)
		} else {
			if err := dbf.Close(); err != nil {
				panic(err)
			}

			return dbf.Name()
		}
	}()
	defer os.Remove(dbpath)

	// open database
	cfg := &Config{
		DatabasePath: dbpath,
	}
	db, err := OpenBoltDatabase(cfg)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	// add mappings
	for _, m := range boltdbMappings {
		if err := db.AddMapping(&m); err != nil {
			panic(err)
		}
	}

	// test mappings
	for _, m := range boltdbMappings {
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
		if len(v) != len(boltdbMappings) {
			t.Errorf("Bad mapping count %v, expected %v", len(v), len(boltdbMappings))
		}
	}

	// delete mappings
	for _, m := range boltdbMappings {
		if err := db.DeleteMapping(m.Key); err != nil {
			panic(err)
		}

		if _, err := db.GetMapping(m.Key); err != MappingNotFoundError {
			t.Errorf("Mapping was deleted but still exists in database")
		}
	}
}
