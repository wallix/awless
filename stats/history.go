package stats

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/boltdb/bolt"
)

type Line struct {
	Id      int
	Command []string
	Time    time.Time
}

var historyBucketName string

func init() {
	historyBucketName = "line"
}

func (db *DB) FlushHistory() error {
	return db.DeleteBucket(historyBucketName)
}

func (db *DB) GetHistory(fromCommandId int) ([]*Line, error) {
	var result []*Line
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(historyBucketName))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.Seek(itob(fromCommandId)); k != nil; k, v = c.Next() {
			line := &Line{}
			e := json.Unmarshal(v, line)
			if e != nil {
				return e
			}
			result = append(result, line)
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (db *DB) AddHistoryCommandWithTime(command []string, time time.Time) error {
	line := Line{Command: command, Time: time}

	err := db.Update(func(tx *bolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(historyBucketName))
		if e != nil {
			return e
		}

		id, e := b.NextSequence()
		if e != nil {
			return e
		}
		line.Id = int(id)

		buf, e := json.Marshal(line)
		if e != nil {
			return e
		}
		return b.Put(itob(line.Id), buf)
	})

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) AddHistoryCommand(command []string) error {
	return db.AddHistoryCommandWithTime(command, time.Now())
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
