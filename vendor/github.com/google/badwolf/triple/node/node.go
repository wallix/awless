// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package node provides the abstraction to build and use BadWolf nodes.
package node

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

// Type describes the type of the node.
type Type string

// String converts a type to its string form.
func (t *Type) String() string {
	return string(*t)
}

// Covariant checks for given two types A and B, A covariant B if B _is a_ A.
// In other word, A _covariant_ B if B is a prefix of A.
func (t *Type) Covariant(ot *Type) bool {
	if !strings.HasPrefix(t.String(), ot.String()) {
		return false
	}
	// /type/foo is covariant of /type, but /typefoo is not covariant of /type.
	return len(t.String()) == len(ot.String()) || t.String()[len(ot.String())] == '/'
}

// ID represents a node ID.
type ID string

// String converts a ID to its string form.
func (i *ID) String() string {
	return string(*i)
}

// Node describes a node in a BadWolf graph.
type Node struct {
	t  *Type
	id *ID
}

// Type returns the type of the node.
func (n *Node) Type() *Type {
	return n.t
}

// ID returns the ID of the node.
func (n *Node) ID() *ID {
	return n.id
}

// String returns a pretty printing representation of Node.
func (n *Node) String() string {
	return fmt.Sprintf("%s<%s>", n.t.String(), n.id.String())
}

// Parse returns a node given a pretty printed representation of Node.
func Parse(s string) (*Node, error) {
	raw := strings.TrimSpace(s)
	idx := strings.Index(raw, "<")
	if idx < 0 {
		return nil, fmt.Errorf("node.Parser: invalid format, could not find ID in %v", raw)
	}
	t, err := NewType(raw[:idx])
	if err != nil {
		return nil, fmt.Errorf("node.Parser: invalid type %q, %v", raw[:idx], err)
	}
	if raw[len(raw)-1] != '>' {
		return nil, fmt.Errorf("node.Parser: pretty printing should finish with '>' in %q", raw)
	}
	id, err := NewID(raw[idx+1 : len(raw)-1])
	if err != nil {
		return nil, fmt.Errorf("node.Parser: invalid ID in %q, %v", raw, err)
	}
	return NewNode(t, id), nil
}

// Covariant checks if the types of two nodes is covariant.
func (n *Node) Covariant(on *Node) bool {
	return n.t.Covariant(on.t)
}

// NewType creates a new type from plain string.
func NewType(t string) (*Type, error) {
	if strings.ContainsAny(t, " \t\n\r") {
		return nil, fmt.Errorf("node.NewType(%q) does not allow spaces", t)
	}
	if !strings.HasPrefix(t, "/") || strings.HasSuffix(t, "/") {
		return nil, fmt.Errorf("node.NewType(%q) should start with a '/' and do not end with '/'", t)
	}
	if t == "" {
		return nil, fmt.Errorf("node.NewType(%q) cannot create empty types", t)
	}
	nt := Type(t)
	return &nt, nil
}

// NewID create a new ID from a plain string.
func NewID(id string) (*ID, error) {
	if strings.ContainsAny(id, "<>") {
		return nil, fmt.Errorf("node.NewID(%q) does not allow '<' or '>'", id)
	}
	if id == "" {
		return nil, fmt.Errorf("node.NewID(%q) cannot create empty ID", id)
	}
	nID := ID(id)
	return &nID, nil
}

// NewNode returns a new node constructed from a type and an ID.
func NewNode(t *Type, id *ID) *Node {
	return &Node{
		t:  t,
		id: id,
	}
}

// NewNodeFromStrings returns a new node constructed from a type and ID
// represented as plain strings.
func NewNodeFromStrings(sT, sID string) (*Node, error) {
	t, err := NewType(sT)
	if err != nil {
		return nil, err
	}
	n, err := NewID(sID)
	if err != nil {
		return nil, err
	}
	return NewNode(t, n), nil
}

const chanSize = 256

// The channel to recover the next unique value used to create a blank node.
var (
	nextVal chan uuid.UUID
	tBlank  Type
)

func init() {
	var buffer bytes.Buffer
	b := make([]byte, 16)

	// Get the current user name.
	osU, err := user.Current()
	user := "UNKNOW"
	if err == nil {
		user = osU.Username
	}
	buffer.WriteString(user)

	// Create the constant to make build a unique ID.
	start := uint64(time.Now().UnixNano())
	binary.PutUvarint(b, start)
	buffer.Write(b)

	pid := uint64(os.Getpid())
	binary.PutUvarint(b, pid)
	buffer.Write(b)

	// Set the node.
	if !uuid.SetNodeID(buffer.Bytes()) {
		os.Exit(-1)
	}

	// Initialize the channel and blank node type.
	nextVal, tBlank = make(chan uuid.UUID, chanSize), Type("/_")

	go func() {
		for {
			nextVal <- uuid.NewUUID()
		}
	}()
}

// NewBlankNode creates a new blank node. The blank node ID is guaranteed to
// be unique in BadWolf.
func NewBlankNode() *Node {
	uuid := <-nextVal
	id := ID(uuid.String())
	return &Node{
		t:  &tBlank,
		id: &id,
	}
}

// UUID returns a global unique identifier for the given node. It is
// implemented as the SHA1 UUID of the node values.
func (n *Node) UUID() uuid.UUID {
	var buffer bytes.Buffer
	buffer.WriteString(string(*n.t))
	buffer.WriteString(string(*n.id))
	return uuid.NewSHA1(uuid.NIL, buffer.Bytes())
}
