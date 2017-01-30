package database

import (
	"encoding/json"
	"fmt"
	"github.com/wallix/awless/template"

	"github.com/boltdb/bolt"
)

const EXECUTIONS_BUCKET = "executions"

func (db *DB) AddTemplateExecution(templ *template.TemplateExecution) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(EXECUTIONS_BUCKET))
		if err != nil {
			return fmt.Errorf("create bucket %s: %s", EXECUTIONS_BUCKET, err)
		}

		b, err := json.Marshal(templ)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(templ.ID), b)
	})
}

func (db *DB) DeleteTemplateExecutions() error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(EXECUTIONS_BUCKET))
	})
}

func (db *DB) GetTemplateExecutions() ([]*template.TemplateExecution, error) {
	var result []*template.TemplateExecution

	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(EXECUTIONS_BUCKET))
		if b == nil {
			return nil
		}

    c := b.Cursor()

    for k, v := c.First(); k != nil; k, v = c.Next() {
    		t := &template.TemplateExecution{}
    		if err := json.Unmarshal(v, t); err != nil {
    			return err
    		}
        result = append(result, t)
    }

		return nil
	})

	return result, err
}

