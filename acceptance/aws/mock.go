package awsat

import (
	"reflect"
	"testing"
)

type mock interface {
	Calls() map[string]int
	SetInputs(map[string]interface{})
	SetIgnored(map[string]struct{})
	SetTesting(*testing.T)
}

type basicMock struct {
	t             *testing.T
	calls         map[string]int
	expInputs     map[string]interface{}
	ignoredInputs map[string]struct{}
}

func (m *basicMock) addCall(call string) {
	if m.calls == nil {
		m.calls = make(map[string]int)
	}
	m.calls[call]++
}

func (m *basicMock) Calls() map[string]int {
	return m.calls
}

func (m *basicMock) SetTesting(t *testing.T) {
	m.t = t
}

func (m *basicMock) SetInputs(inputs map[string]interface{}) {
	m.expInputs = inputs
}

func (m *basicMock) SetIgnored(ignored map[string]struct{}) {
	m.ignoredInputs = ignored
}

func (m *basicMock) verifyInput(call string, got interface{}) {
	m.t.Helper()
	if m.expInputs == nil {
		return
	}
	if _, isIgnored := m.ignoredInputs[call]; isIgnored {
		return
	}
	if want := m.expInputs[call]; !reflect.DeepEqual(want, got) {
		m.t.Fatalf("got %#v, want %#v", got, want)
	}
}
