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

func (db *DB) FlushLogs() error {
	return db.SetValue(logsKey, []byte{})
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
	return db.SetValue(logsKey, b)
}

func (db *DB) GetLogs() (logs []*Log, err error) {
	b, err := db.GetValue(logsKey)
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
