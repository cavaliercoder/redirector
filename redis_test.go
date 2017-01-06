package main

import (
	"testing"
)

func TestRedis(t *testing.T) {
	// open database
	cfg := &Config{
		DatabasePath: "localhost:6379",
	}
	db, err := OpenRedisDatabase(cfg)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}()

	testDB(t, db)
}
