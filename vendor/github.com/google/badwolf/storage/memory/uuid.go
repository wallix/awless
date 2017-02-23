package memory

import (
	"encoding/base64"
	"fmt"

	"github.com/pborman/uuid"
)

// UUIDToBase64 recodes it into a compact string.
func UUIDToBase64(uuid uuid.UUID) string {
	return base64.StdEncoding.EncodeToString([]byte(uuid))
}

// Base64ToUUID attempts to decode a base 64 encoded UUID.
func Base64ToUUID(b64 string) (uuid.UUID, error) {
	bs, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	if got, want := len(bs), 16; got != want {
		return nil, fmt.Errorf("invalid UUID size; got %d, want %d", got, want)
	}
	return uuid.UUID(bs), nil
}

// UUIDToByteString converts the bytes of the UUID into a, likely,
// unreadable string. This is useful to minimized the key size used
// for the inmemory indices.
func UUIDToByteString(uuid uuid.UUID) string {
	return string(uuid)
}
