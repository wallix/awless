package stats

import (
	"encoding/json"
	"time"
)

const (
	LOGS_KEY = "logs"
)

type Log struct {
	Msg  string
	Hits int
	Date time.Time
}

func (db *DB) FlushLogs() error {
	return db.SetValue(LOGS_KEY, []byte{})
}

func (db *DB) AddLog(msg string) error {
	logs, err := db.GetLogs()
	if err != nil {
		return err
	}
	log := findLogInSlice(msg, logs)
	if log == nil {
		log = &Log{Msg: msg, Date: time.Now()}
		logs = append(logs, log)
	}
	log.Hits++

	b, err := json.Marshal(logs)
	if err != nil {
		return err
	}
	return db.SetValue(LOGS_KEY, b)
}

func (db *DB) GetLogs() (logs []*Log, err error) {
	b, err := db.GetValue(LOGS_KEY)
	if err != nil {
		return logs, err
	}
	if len(b) == 0 {
		return logs, nil
	}
	err = json.Unmarshal(b, &logs)
	if err != nil {
		return logs, err
	}

	return logs, err
}

func findLogInSlice(msg string, logs []*Log) *Log {
	for _, log := range logs {
		if log.Msg == msg {
			return log
		}
	}
	return nil
}
