package display

import "fmt"

type Header interface {
	propKey() string
	title() string
	format(i interface{}) string
}

type StringHeader struct {
	Prop, Friendly string
}

func (h StringHeader) format(i interface{}) string {
	if i == nil {
		return ""
	}

	return fmt.Sprint(i)
}
func (h StringHeader) propKey() string { return h.Prop }
func (h StringHeader) title() string {
	if h.Friendly == "" {
		return h.Prop
	}
	return h.Friendly
}
