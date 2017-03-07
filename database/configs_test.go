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
	"reflect"
	"testing"
)

func TestLoadConfigs(t *testing.T) {
	db, close := newTestDb()
	defer close()
	configKey := "config"
	d, err := db.GetConfigs(configKey)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := d, make(configs); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	db.SetConfig(configKey, "key-1", "value-1")
	db.SetConfig(configKey, "key-2", "value-2")
	db.SetConfig(configKey, "key-1", "value-3")

	expected := configs{
		"key-1": "value-3",
		"key-2": "value-2",
	}

	d, err = db.GetConfigs(configKey)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := d, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	v, _ := db.GetConfig(configKey, "key-1")
	if got, want := v.(string), "value-3"; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	v, _ = db.GetConfig(configKey, "key-2")
	if got, want := v.(string), "value-2"; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	db.UnsetConfig(configKey, "key-2")

	expected = configs{
		"key-1": "value-3",
	}

	d, err = db.GetConfigs(configKey)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := d, expected; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	v, ok := db.GetConfig(configKey, "key-1")
	if got, want := v.(string), "value-3"; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if got, want := ok, true; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	str, ok := db.GetConfigString(configKey, "key-1")
	if got, want := str, "value-3"; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if got, want := ok, true; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	_, ok = db.GetConfig(configKey, "key-2")
	if got, want := ok, false; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	_, ok = db.GetConfigString(configKey, "key-2")
	if got, want := ok, false; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
