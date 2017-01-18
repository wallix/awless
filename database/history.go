package database

import (
	"encoding/binary"
	"time"
)

// DeleteHistory empties the history from database
func (db *DB) DeleteHistory() error {
	return db.DeleteBucket(historyBucketName)
}

// GetHistory gets the history from database
func (db *DB) GetHistory(fromID int) ([]*line, error) {
	return db.getLinesFromBucket(historyBucketName, fromID)
}

// AddHistoryCommand adds a command to history in database
func (db *DB) AddHistoryCommand(command []string) error {
	return db.AddHistoryCommandWithTime(command, time.Now())
}

// AddHistoryCommandWithTime adds a command to history in database where time can be set
func (db *DB) AddHistoryCommandWithTime(command []string, time time.Time) error {
	l := line{Command: command, Time: time}

	return db.addLineToBucket(historyBucketName, l)
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
