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
