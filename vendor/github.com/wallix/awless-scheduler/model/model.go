package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/adler32"
	"time"
)

const (
	AwlessFileExt = "aws"
	StampLayout   = "2006-01-02-15h04m05s"
)

type ServiceInfo struct {
	Uptime          string
	ServiceAddr     string
	TickerFrequency string
	UnixSockMode    bool
}

type Task struct {
	Content  string
	RunAt    time.Time
	RevertAt time.Time
	Region   string
}

func (tk *Task) AsFilename() string {
	checksum := adler32.Checksum([]byte(tk.Content))
	return fmt.Sprintf("%d_%s_%s_%s.%s", checksum, tk.RunAt.UTC().Format(StampLayout), tk.RevertAt.UTC().Format(StampLayout), tk.Region, AwlessFileExt)
}

func (tk *Task) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	jsonValue, err := json.Marshal(tk.Content)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"Content\":%s,", jsonValue))
	if !tk.RunAt.IsZero() {
		jsonValue, err = json.Marshal(tk.RunAt)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"RunAt\":%s,", jsonValue))
		buffer.WriteString(fmt.Sprintf("\"RunIn\":\"%s\",", time.Until(tk.RunAt)))
	}
	if !tk.RevertAt.IsZero() {
		jsonValue, err = json.Marshal(tk.RevertAt)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"RevertAt\":%s,", jsonValue))
		buffer.WriteString(fmt.Sprintf("\"RevertIn\":\"%s\",", time.Until(tk.RevertAt)))
	}
	buffer.WriteString(fmt.Sprintf("\"Region\":\"%s\"", tk.Region))

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
