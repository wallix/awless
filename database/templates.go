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
	"errors"
	"fmt"

	"github.com/wallix/awless/template"

	"github.com/boltdb/bolt"
)

const TEMPLATES_BUCKET = "templates"

func (db *DB) AddTemplate(tplExec *template.TemplateExecution) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		if tplExec.ID == "" {
			return errors.New("cannot persist template with empty ID")
		}

		bucket, err := tx.CreateBucketIfNotExists([]byte(TEMPLATES_BUCKET))
		if err != nil {
			return fmt.Errorf("create bucket %s: %s", TEMPLATES_BUCKET, err)
		}

		b, err := tplExec.MarshalJSON()
		if err != nil {
			return err
		}

		return bucket.Put([]byte(tplExec.ID), b)
	})
}

func (db *DB) GetTemplate(id string) (*template.TemplateExecution, error) {
	tplExec := &template.TemplateExecution{}

	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TEMPLATES_BUCKET))
		if b == nil {
			return errors.New("no templates stored yet")
		}
		if content := b.Get([]byte(id)); content != nil {
			return tplExec.UnmarshalJSON(content)
		} else {
			return fmt.Errorf("no content for id '%s'", id)
		}
	})

	return tplExec, err
}

func (db *DB) DeleteTemplates() error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TEMPLATES_BUCKET))
		if b == nil {
			return nil
		}
		return tx.DeleteBucket([]byte(TEMPLATES_BUCKET))
	})
}

func (db *DB) DeleteTemplate(id string) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TEMPLATES_BUCKET))
		if b == nil {
			return errors.New("no templates stored yet")
		}
		return b.Delete([]byte(id))
	})
}

type LoadedTemplate struct {
	Err      error
	TplExec  *template.TemplateExecution
	Key, Raw string
}

func (db *DB) ListTemplates() ([]*LoadedTemplate, error) {
	var results []*LoadedTemplate

	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TEMPLATES_BUCKET))
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			tplExec := &template.TemplateExecution{}
			terr := tplExec.UnmarshalJSON(v)
			lt := &LoadedTemplate{TplExec: tplExec, Err: terr, Key: string(k), Raw: string(v)}
			results = append(results, lt)
		}

		return nil
	})

	return results, err
}
