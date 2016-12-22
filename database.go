package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"time"
)

var (
	MappingNotFoundError = fmt.Errorf("Mapping not found")
)

var (
	MAPPINGS_BUCKET = []byte("mappings")
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
		_, err := tx.CreateBucketIfNotExists([]byte("mappings"))
		if err != nil {
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
func (db *Database) get(b, k []byte, v interface{}) error {
	return db.bdb.View(func(tx *bolt.Tx) error {
		vb := tx.Bucket(b).Get(k)
		if vb == nil {
			return MappingNotFoundError
		}

		return UnmarshallBinary(vb, v)
	})
}

func (db *Database) add(b, k []byte, v interface{}) error {
	return db.bdb.Update(func(tx *bolt.Tx) error {
		vb, err := MarshallBinary(v)
		if err != nil {
			return err
		}

		return tx.Bucket(b).Put(k, vb)
	})
}

func (db *Database) AddMapping(m *Mapping) error {
	return db.add(MAPPINGS_BUCKET, []byte(m.Key), m)
}

func (db *Database) GetMapping(key string) (*Mapping, error) {
	m := &Mapping{}
	if err := db.get(MAPPINGS_BUCKET, []byte(key), m); err != nil {
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
