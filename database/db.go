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
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

const (
	Filename     = "awless.db"
	awlessBucket = "awless"
)

type DB struct {
	bolt *bolt.DB
}

func Execute(fn func(*DB) error) error {
	db, err := current()
	if err != nil {
		return err
	}
	defer db.Close()

	return fn(db)
}

func current() (*DB, error) {
	awlessHome := os.Getenv("__AWLESS_HOME")
	if awlessHome == "" {
		return nil, errors.New("database: awless home is not set")
	}

	path := filepath.Join(awlessHome, Filename)
	db, err := open(path)
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("db is nil while no error in opening at '%s'", path)
	}

	return db, nil
}

func open(path string) (*DB, error) {
	boltdb, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("opening db at %s: %s (any awless existing process running?)", path, err)
	}

	return &DB{bolt: boltdb}, nil
}

// DeleteBucket deletes a bucket if it exists
func (db *DB) DeleteBucket(name string) error {
	return db.deleteBucket(name)
}

// GetBytes gets a []byte value from database
func (db *DB) GetBytes(key string) ([]byte, error) {
	return db.getValue(key)
}

// GetStringValue gets a string value from database
func (db *DB) GetStringValue(key string) (string, error) {
	str, err := db.getValue(key)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

// GetTimeValue gets a time value from database
func (db *DB) GetTimeValue(key string) (time.Time, error) {
	var t time.Time
	bin, err := db.getValue(key)
	if err != nil {
		return t, err
	}
	if len(bin) == 0 {
		return t, nil
	}
	err = t.UnmarshalBinary(bin)
	return t, err
}

// GetIntValue gets a int value from database
func (db *DB) GetIntValue(key string) (int, error) {
	str, err := db.GetStringValue(key)
	if err != nil {
		return 0, err
	}
	if str == "" {
		return 0, nil
	}
	return strconv.Atoi(str)
}

// SetBytes sets a []byte value in database
func (db *DB) SetBytes(key string, value []byte) error {
	return db.setValue(key, value)
}

// SetStringValue sets a string value in database
func (db *DB) SetStringValue(key, value string) error {
	return db.setValue(key, []byte(value))
}

// SetTimeValue sets a time value in database
func (db *DB) SetTimeValue(key string, t time.Time) error {
	bin, err := t.MarshalBinary()
	if err != nil {
		return err
	}
	return db.setValue(key, bin)
}

// SetIntValue sets a int value in database
func (db *DB) SetIntValue(key string, value int) error {
	return db.SetStringValue(key, strconv.Itoa(value))
}

// Close the database
func (db *DB) Close() {
	if db.bolt != nil {
		db.bolt.Close()
	}
}

func (db *DB) deleteBucket(name string) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(name))
		if b == nil {
			return nil
		}
		e := tx.DeleteBucket([]byte(name))
		return e
	})
}

func (db *DB) getValue(key string) ([]byte, error) {
	var value []byte
	err := db.bolt.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket([]byte(awlessBucket)); b != nil {
			value = b.Get([]byte(key))
		}
		return nil
	})
	if err != nil {
		return value, err
	}

	return value, nil
}

func (db *DB) setValue(key string, value []byte) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(awlessBucket))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
}
