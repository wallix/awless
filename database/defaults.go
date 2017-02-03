package database

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

const (
	SyncAuto         = "sync.auto"
	RegionKey        = "region"
	InstanceTypeKey  = "instance.type"
	InstanceImageKey = "instance.image"
	InstanceCountKey = "instance.count"
)

type defaults map[string]interface{}

func MustGetDefaultRegion() string {
	db, close := Current()
	defer close()
	return db.MustGetDefaultRegion()
}

func (db *DB) MustGetDefaultRegion() string {
	region, ok := db.GetDefaultString(RegionKey)
	if !ok {
		fmt.Fprintf(os.Stderr, "config: missing region. Set it with `awless config set region`\n")
		os.Exit(-1)
	}
	return region
}

func (db *DB) GetDefaults() (defaults, error) {
	d := make(defaults)
	b, err := db.GetBytes(defaultsKey)
	if err != nil {
		return d, err
	}
	if len(b) == 0 {
		return d, nil
	}

	dec := gob.NewDecoder(bytes.NewReader(b))
	if err := dec.Decode(&d); err != nil {
		return d, nil
	}
	return d, err
}

func (db *DB) SetDefault(k string, v interface{}) error {
	d, err := db.GetDefaults()
	if err != nil {
		return err
	}
	d[k] = v
	return db.saveDefaults(d)
}

func (db *DB) UnsetDefault(k string) error {
	d, err := db.GetDefaults()
	if err != nil {
		return err
	}
	delete(d, k)
	return db.saveDefaults(d)
}

func (db *DB) GetDefault(k string) (interface{}, bool) {
	d, err := db.GetDefaults()
	if err != nil {
		return nil, false
	}
	i, ok := d[k]
	return i, ok
}

func (db *DB) GetDefaultString(k string) (string, bool) {
	v, ok := db.GetDefault(k)
	if !ok {
		return "", ok
	}
	str, ok := v.(string)
	return str, ok
}

func (db *DB) saveDefaults(d defaults) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(d); err != nil {
		return err
	}
	return db.SetBytes(defaultsKey, buff.Bytes())
}
