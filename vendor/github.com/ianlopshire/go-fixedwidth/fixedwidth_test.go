package fixedwidth

var (
	nilFloat64 *float64
	nilFloat32 *float32
	nilInt     *int
	nilString  *string
)

func float64p(v float64) *float64 { return &v }
func float32p(v float32) *float32 { return &v }
func intp(v int) *int             { return &v }
func stringp(v string) *string    { return &v }

// EncodableString is a string that implements the encoding TextUnmarshaler and TextMarshaler interface.
// This is useful for testing.
type EncodableString struct {
	S   string
	Err error
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *EncodableString) UnmarshalText(text []byte) error {
	s.S = string(text)
	return s.Err
}

// MarshalText implements encoding.TextUnmarshaler.
func (s EncodableString) MarshalText() ([]byte, error) {
	return []byte(s.S), s.Err
}
