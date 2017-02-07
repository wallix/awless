package database

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/wallix/awless/cloud/aws"
)

const (
	salt             = "bg6B8yTTq8chwkN0BqWnEzlP4OkpcQDhO45jUOuXm1zsNGDLj3"
	databaseFilename = "awless.db"
)

// A DB stores awless config, logs...
type DB struct {
	bolt *bolt.DB
}

func Current() (*DB, func()) {
	awlessHome := os.Getenv("__AWLESS_HOME")
	if awlessHome == "" {
		fmt.Fprintf(os.Stderr, "database: awless home is not set\n")
		os.Exit(-1)
	}
	db, err := open(filepath.Join(awlessHome, databaseFilename))
	if err != nil {
		fmt.Fprintf(os.Stderr, "database: %s\n", err)
		os.Exit(-1)
	}
	todefer := func() {
		db.Close()
	}
	return db, todefer
}

func open(path string) (*DB, error) {
	boltdb, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("opening db at %s: %s (any awless existing process running?)", path, err)
	}

	return &DB{bolt: boltdb}, nil
}

func InitDB(firstInstall bool) error {
	db, closing := Current()
	defer closing()
	if db == nil {
		return fmt.Errorf("database: empty current database")
	}
	id, err := db.GetStringValue(AwlessIdKey)
	if err != nil || id == "" {
		userID, err := aws.SecuAPI.GetUserId()
		if err != nil {
			return err
		}
		newID, err := generateAnonymousID(userID)
		if err != nil {
			return err
		}
		if err = db.SetStringValue(AwlessIdKey, newID); err != nil {
			return err
		}
		accountID, err := aws.SecuAPI.GetAccountId()
		if err != nil {
			return err
		}
		aID, err := generateAnonymousID(accountID)
		if err != nil {
			return err
		}
		if err = db.SetStringValue(AwlessAIdKey, aID); err != nil {
			return err
		}
	}

	return nil
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

func (db *DB) addLineToBucket(bucket string, l line) error {
	return db.bolt.Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(bucket))
		if e != nil {
			return e
		}

		id, e := b.NextSequence()
		if e != nil {
			return e
		}
		l.ID = int(id)

		buf, e := json.Marshal(l)
		if e != nil {
			return e
		}
		return b.Put(itob(l.ID), buf)
	})
}

func (db *DB) getLinesFromBucket(bucket string, fromID int) ([]*line, error) {
	var result []*line
	err := db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(historyBucketName))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.Seek(itob(fromID)); k != nil; k, v = c.Next() {
			l := &line{}
			e := json.Unmarshal(v, l)
			if e != nil {
				return e
			}
			result = append(result, l)
		}
		return nil
	})
	return result, err
}

func generateAnonymousID(seed string) (string, error) {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(salt+seed))), nil
}
