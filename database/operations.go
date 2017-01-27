package database

import (
	"encoding/json"
	"fmt"
	"github.com/wallix/awless/template"

	"github.com/boltdb/bolt"
)

var operationsBucket = "operations"

func (db *DB) AddTemplateOperation(templ *template.Template) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(operationsBucket))
		if err != nil {
			return fmt.Errorf("create bucket %s: %s", operationsBucket, err)
		}

		b, err := json.Marshal(templ)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(templ.ID), b)
	})
}

func (db *DB) GetTemplateOperations() ([]string, error) {
	var result []string
	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(operationsBucket))
		if b == nil {
			return nil
		}

    c := b.Cursor()

    for k, v := c.First(); k != nil; k, v = c.Next() {
        result = append(result, string(v))
    }

		return nil
	})

	return result, err
}

