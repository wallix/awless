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
	"time"
)

// A Log represents a log of an error that occured in awless
type Log struct {
	Msg  string
	Hits int
	Date time.Time
}

func (db *DB) DeleteLogs() error {
	return db.SetBytes(logsKey, []byte{})
}

func (db *DB) AddLog(msg string) error {
	logs, err := db.GetLogs()
	if err != nil {
		return err
	}
	l := findLogInSlice(msg, logs)
	if l == nil {
		l = &Log{Msg: msg, Date: time.Now()}
		logs = append(logs, l)
	}
	l.Hits++

	b, err := json.Marshal(logs)
	if err != nil {
		return err
	}
	return db.SetBytes(logsKey, b)
}

func (db *DB) GetLogs() (logs []*Log, err error) {
	b, err := db.GetBytes(logsKey)
	if err != nil {
		return logs, err
	}
	if len(b) == 0 {
		return logs, nil
	}
	err = json.Unmarshal(b, &logs)
	return logs, err
}

func findLogInSlice(msg string, logs []*Log) *Log {
	for _, l := range logs {
		if l.Msg == msg {
			return l
		}
	}
	return nil
}
