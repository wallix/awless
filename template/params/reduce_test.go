package params

import (
	"reflect"
	"testing"
)

func TestReduce(t *testing.T) {
	data := map[string]interface{}{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	plusOne := func(in map[string]interface{}) (out map[string]interface{}, err error) {
		out = make(map[string]interface{})
		for k, i := range in {
			out[k] = i.(int) + 1
		}
		return
	}
	red := newReducer(plusOne, "one", "three")
	out, err := red.Reduce(data)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := red.Keys(), []string{"one", "three"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := len(out), 2; got != want {
		t.Fatalf("got length %d, want length %d", got, want)
	}
	if v, ok := out["one"].(int); !ok || v != 2 {
		t.Fatalf("invalid content: %v", out)
	}
	if v, ok := out["three"].(int); !ok || v != 4 {
		t.Fatalf("invalid content: %v", out)
	}
}
