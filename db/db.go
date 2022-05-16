package db

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("default")

type Database struct {
	db *bolt.DB
}

func NewDatabase(dbPath string) (db *Database, closeFunc func() error, err error) {
	boltDB, err := bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, nil, err
	}
	db = &Database{
		db: boltDB,
	}
	closeFunc = boltDB.Close

	if err := db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("creating default bucket: %w", err)
	}

	return db, closeFunc, nil
}

func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		return b.Put([]byte(key), value)
	})
}

func (d *Database) GetKey(key string) ([]byte, error) {
	var result []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(defaultBucket)
		result = b.Get([]byte(key))
		return nil
	})

	if err == nil {
		return result, err
	}

	return nil, err
}
