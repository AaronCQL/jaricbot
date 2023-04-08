package database

import (
	"encoding/binary"
	"path/filepath"

	"github.com/cockroachdb/pebble"
	"github.com/vmihailenco/msgpack/v5"
)

type Database struct {
	db *pebble.DB
}

func New(dataDir, store string) *Database {
	db, err := pebble.Open(filepath.Join(dataDir, store), &pebble.Options{})
	if err != nil {
		panic(err)
	}
	return &Database{
		db: db,
	}
}

func int64ToBytes(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

func (d *Database) Get(key int64, value any) error {
	data, closer, err := d.db.Get(int64ToBytes(key))
	if err != nil {
		return err
	}
	defer closer.Close()
	return msgpack.Unmarshal(data, value)
}

func (d *Database) Set(key int64, value any) error {
	bytes, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	return d.db.Set(int64ToBytes(key), bytes, pebble.Sync)
}

func (d *Database) Close() error {
	return d.db.Close()
}
