package stats

import (
	"time"

	"github.com/boltdb/bolt"
)

const AWLESS_BUCKET = "awless"

type DB struct {
	*bolt.DB
}

func OpenDB(name string) (*DB, error) {
	boltdb, err := bolt.Open(name, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	db := &DB{boltdb}

	if id, err := db.GetStringValue(AWLESS_ID_KEY); err != nil {
		return nil, err
	} else if id == "" {
		newId, err := generateAwlessId()
		if err != nil {
			return nil, err
		}
		if err = db.SetStringValue(AWLESS_ID_KEY, newId); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (db *DB) DeleteBucket(name string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		if b == nil {
			return nil
		}
		e := tx.DeleteBucket([]byte(name))
		return e
	})
}

func (db *DB) GetStringValue(key string) (string, error) {
	var value []byte
	err := db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(AWLESS_BUCKET)); b != nil {
			value = b.Get([]byte(key))
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return string(value), nil
}

func (db *DB) SetStringValue(key, value string) error {
	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(AWLESS_BUCKET))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
}
