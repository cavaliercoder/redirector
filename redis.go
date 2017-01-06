package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
)

type RedisDatabase struct {
	client redis.Conn
}

func OpenRedisDatabase(cfg *Config) (Database, error) {
	client, err := redis.Dial("tcp", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	return &RedisDatabase{
		client: client,
	}, nil
}

func (db *RedisDatabase) Close() error {
	return db.client.Close()
}

func (db *RedisDatabase) Stats() interface{} {
	s, err := redis.String(db.client.Do("INFO", "all"))
	if err != nil {
		return nil
	}

	stats := make(map[string]string, 0)
	lines := strings.Split(s, "\r\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") {
			v := strings.SplitN(line, ":", 2)
			if len(v) == 2 {
				stats[v[0]] = v[1]
			}
		}
	}

	return stats
}

func (db *RedisDatabase) Count() int {
	str, err := redis.String(db.client.Do("INFO", "keyspace"))
	if err != nil {
		panic(err)
	}

	i := -1
	o := strings.Index(str, ":keys=") + 6
	fmt.Sscanf(str[o:], "%d", &i)

	return i
}

func (db *RedisDatabase) AddMapping(m *Mapping) error {
	b, err := MarshallBinary(m)
	if err != nil {
		return err
	}

	if res, err := redis.String(db.client.Do("SET", m.Key, b)); err != nil {
		return err
	} else if res != "OK" {
		return fmt.Errorf("Failed to set key %v: %v", m.Key, res)
	}

	return nil
}

func (db *RedisDatabase) GetMapping(key string) (*Mapping, error) {
	b, err := redis.Bytes(db.client.Do("GET", key))
	if err == redis.ErrNil {
		return nil, MappingNotFoundError
	}
	if err != nil {
		return nil, err
	}

	m := &Mapping{}
	if err := UnmarshallBinary(b, m); err != nil {
		return nil, err
	}

	return m, nil
}

func (db *RedisDatabase) GetMappings() ([]*Mapping, error) {
	values, err := redis.Values(db.client.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}

	var keys = make([]string, 0)
	if err := redis.ScanSlice(values, &keys); err != nil {
		return nil, err
	}

	// TODO: improve performance of RedisDatabase.GetMappings
	mappings := make([]*Mapping, len(keys))
	for i, key := range keys {
		b, err := redis.Bytes(db.client.Do("GET", key))
		if err != nil {
			return nil, err
		}

		m := &Mapping{}
		if err := UnmarshallBinary(b, m); err != nil {
			return nil, err
		}

		mappings[i] = m
	}

	return mappings, nil
}

func (db *RedisDatabase) DeleteMapping(key string) error {
	i, err := redis.Int(db.client.Do("DEL", key))
	if err != nil {
		return err
	}

	if i != 1 {
		return fmt.Errorf("Expected to delete 1 mapping, got %v", i)
	}

	return nil
}
