package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

var (
	MappingNotFoundError = fmt.Errorf("Mapping not found")
)

type Database struct {
	bdb *bolt.DB
}

func OpenDatabase(cfg *Config) (*Database, error) {
	options := &bolt.Options{
		Timeout: 3 * time.Second,
	}

	bdb, err := bolt.Open(cfg.DatabasePath, 0600, options)
	if err != nil {
		return nil, err
	}

	bdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("mappings"))
		if err != nil && err != bolt.ErrBucketExists {
			return err
		}

		return nil
	})

	return &Database{
		bdb: bdb,
	}, nil
}

func (db *Database) Close() error {
	return db.bdb.Close()
}

func (db *Database) Stats() interface{} {
	return db.bdb.Stats()
}

func (db *Database) AddMapping(m *Mapping) error {
	return db.bdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("mappings"))
		return b.Put([]byte(m.Key), []byte(m.Destination))
	})
}

func (db *Database) GetMapping(key string) (*Mapping, error) {
	var m *Mapping = nil
	if err := db.bdb.View(func(tx *bolt.Tx) error {
		v := tx.Bucket([]byte("mappings")).Get([]byte(key))
		if len(v) == 0 {
			return MappingNotFoundError
		}

		m = &Mapping{
			Key:         key,
			Destination: string(v),
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return m, nil
}

func (db *Database) GetMappings() ([]*Mapping, error) {
	mappings := make([]*Mapping, 0)
	if err := db.bdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("mappings"))
		return b.ForEach(func(k, v []byte) error {
			m := &Mapping{
				Key:         string(k),
				Destination: string(v),
			}
			mappings = append(mappings, m)

			return nil
		})
	}); err != nil {
		return nil, err
	}

	return mappings, nil
}
