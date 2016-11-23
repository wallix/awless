package stats

import "github.com/boltdb/bolt"

type DB struct {
	*bolt.DB
}

func NewDB(name string) (*DB, error) {
	boltdb, err := bolt.Open(name, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &DB{boltdb}, nil
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
