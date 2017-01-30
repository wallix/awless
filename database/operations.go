package database

import (
	"encoding/json"
	"fmt"
	"github.com/wallix/awless/template"

	"github.com/boltdb/bolt"
)

const OPERATIONS_BUCKET = "operations"

func (db *DB) AddTemplateOperation(templ *template.Template) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(OPERATIONS_BUCKET))
		if err != nil {
			return fmt.Errorf("create bucket %s: %s", OPERATIONS_BUCKET, err)
		}

		b, err := json.Marshal(templ)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(templ.ID), b)
	})
}

func (db *DB) DeleteTemplateOperations() error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(OPERATIONS_BUCKET))
	})
}

func (db *DB) GetTemplateOperations() ([]*template.Template, error) {
	var result []*template.Template

	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(OPERATIONS_BUCKET))
		if b == nil {
			return nil
		}

    c := b.Cursor()

    for k, v := c.First(); k != nil; k, v = c.Next() {
    		t := &template.Template{}
    		if err := json.Unmarshal(v, t); err != nil {
    			return err
    		}
        result = append(result, t)
    }

		return nil
	})

	return result, err
}

