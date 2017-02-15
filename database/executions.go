/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package database

import (
	"encoding/json"
	"errors"
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

func (db *DB) GetTemplateExecution(id string) (*template.TemplateExecution, error) {
	tpl := &template.TemplateExecution{}

	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(EXECUTIONS_BUCKET))
		if b == nil {
			return errors.New("no template executions stored yet")
		}
		if content := b.Get([]byte(id)); content != nil {
			return json.Unmarshal(b.Get([]byte(id)), tpl)
		} else {
			return fmt.Errorf("no content for id '%s'", id)
		}

		return nil
	})

	return tpl, err
}

func (db *DB) ListTemplateExecutions() ([]*template.TemplateExecution, error) {
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
