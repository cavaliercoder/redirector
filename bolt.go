package main

import (
	"github.com/boltdb/bolt"
	"time"
)

var (
	MAPPINGS_BUCKET = []byte("mappings")
)

// BoltDatabase implements Database to enable storage of URL mappings in a
// memory-mapped BoltDB data store.
type BoltDatabase struct {
	bdb *bolt.DB
}

func OpenBoltDatabase(cfg *Config) (Database, error) {
	options := &bolt.Options{
		Timeout: 3 * time.Second,
	}

	bdb, err := bolt.Open(cfg.DatabasePath, 0600, options)
	if err != nil {
		return nil, err
	}

	bdb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(MAPPINGS_BUCKET)
		if err != nil {
			return err
		}

		return nil
	})

	return &BoltDatabase{
		bdb: bdb,
	}, nil
}

func (db *BoltDatabase) Close() error {
	return db.bdb.Close()
}

func (db *BoltDatabase) Stats() interface{} {
	return db.bdb.Stats()
}

func (db *BoltDatabase) Count() int64 {
	var count int64 = 0
	db.bdb.View(func(tx *bolt.Tx) error {
		stats := tx.Bucket(MAPPINGS_BUCKET).Stats()
		count = int64(stats.KeyN)
		return nil
	})

	return count
}

func (db *BoltDatabase) get(b, k []byte, v interface{}) error {
	return db.bdb.View(func(tx *bolt.Tx) error {
		vb := tx.Bucket(b).Get(k)
		if vb == nil {
			return MappingNotFoundError
		}

		return UnmarshallBinary(vb, v)
	})
}

func (db *BoltDatabase) add(b, k []byte, v interface{}) error {
	return db.bdb.Update(func(tx *bolt.Tx) error {
		vb, err := MarshallBinary(v)
		if err != nil {
			return err
		}

		return tx.Bucket(b).Put(k, vb)
	})
}

func (db *BoltDatabase) delete(b, k []byte) error {
	return db.bdb.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b).Delete(k)
	})
}

func (db *BoltDatabase) AddMapping(m *Mapping) error {
	return db.add(MAPPINGS_BUCKET, []byte(m.Key), m)
}

func (db *BoltDatabase) GetMapping(key string) (*Mapping, error) {
	m := &Mapping{}
	if err := db.get(MAPPINGS_BUCKET, []byte(key), m); err != nil {
		return nil, err
	}

	return m, nil
}

func (db *BoltDatabase) GetMappings() ([]*Mapping, error) {
	mappings := make([]*Mapping, 0)
	if err := db.bdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(MAPPINGS_BUCKET)
		return b.ForEach(func(k, v []byte) error {
			m := &Mapping{}
			if err := UnmarshallBinary(v, m); err != nil {
				return err
			}

			mappings = append(mappings, m)

			return nil
		})
	}); err != nil {
		return nil, err
	}

	return mappings, nil
}

func (db *BoltDatabase) DeleteMapping(key string) error {
	return db.delete(MAPPINGS_BUCKET, []byte(key))
}

func (db *BoltDatabase) DeleteMappings() (int64, error) {
	count := db.Count()
	err := db.bdb.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket(MAPPINGS_BUCKET).DeleteBucket(MAPPINGS_BUCKET); err != nil {
			return err
		}

		if _, err := tx.CreateBucket(MAPPINGS_BUCKET); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, nil
	}

	return count, nil
}
