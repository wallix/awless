package database

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

// EmptyHistory empties the history from database
func (db *DB) EmptyHistory() error {
	return db.DeleteBucket(historyBucketName)
}

// GetHistory gets the history from database
func (db *DB) GetHistory(fromID int) ([]*line, error) {
	var result []*line
	err := db.View(func(tx *bolt.Tx) error {
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
	if err != nil {
		return result, err
	}
	return result, nil
}

// AddHistoryCommand adds a command to history in database
func (db *DB) AddHistoryCommand(command []string) error {
	return db.AddHistoryCommandWithTime(command, time.Now())
}

// AddHistoryCommandWithTime adds a command to history in database where time can be set
func (db *DB) AddHistoryCommandWithTime(command []string, time time.Time) error {
	l := line{Command: command, Time: time}

	err := db.Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(historyBucketName))
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

	if err != nil {
		return err
	}
	return nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

type line struct {
	ID      int
	Command []string
	Time    time.Time
}
