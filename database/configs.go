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
	"bytes"
	"encoding/gob"
)

type configs map[string]interface{}

func (db *DB) GetConfigs(key string) (configs, error) {
	d := make(configs)
	b, err := db.GetBytes(key)
	if err != nil {
		return d, err
	}
	if len(b) == 0 {
		return d, nil
	}

	return d, gob.NewDecoder(bytes.NewReader(b)).Decode(&d)
}

func (db *DB) SetConfig(configsKey, k string, v interface{}) error {
	d, err := db.GetConfigs(configsKey)
	if err != nil {
		return err
	}
	d[k] = v
	return db.saveConfigs(configsKey, d)
}

func (db *DB) UnsetConfig(configsKey, k string) error {
	d, err := db.GetConfigs(configsKey)
	if err != nil {
		return err
	}
	delete(d, k)
	return db.saveConfigs(configsKey, d)
}

func (db *DB) GetConfig(configsKey, k string) (interface{}, bool) {
	d, err := db.GetConfigs(configsKey)
	if err != nil {
		return nil, false
	}
	i, ok := d[k]
	return i, ok
}

func (db *DB) GetConfigString(configsKey, k string) (string, bool) {
	v, ok := db.GetConfig(configsKey, k)
	if !ok {
		return "", ok
	}
	str, ok := v.(string)
	return str, ok
}

func (db *DB) saveConfigs(configsKey string, d configs) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err := enc.Encode(d); err != nil {
		return err
	}
	return db.SetBytes(configsKey, buff.Bytes())
}
