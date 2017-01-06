package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func tmpBoltDB(fn func(Database)) {
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

	fn(db)
}

func TestBoltDB(t *testing.T) {
	tmpBoltDB(func(db Database) {
		testDB(t, db)
	})
}
