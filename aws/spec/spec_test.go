package awsspec

import (
	"strings"
	"testing"
)

func TestEnumValidator(t *testing.T) {
	tcases := []struct {
		validator      *enumValidator
		value          *string
		expErrContains []string
	}{
		{NewEnumValidator("test1"), String("test1"), nil},
		{NewEnumValidator("test1"), String("test2"), []string{"test1", "test2"}},
		{NewEnumValidator("test1", "test2"), String("test1"), nil},
		{NewEnumValidator("test1", "test2"), String("TesT2"), nil},
		{NewEnumValidator("test1", "test2"), String("test3"), []string{"test1", "test2", "test3"}},
		{NewEnumValidator("test1", "test2", "test4"), String("test3"), []string{"test1", "test2", "test3", "test4"}},
	}
	for i, tcase := range tcases {
		err := tcase.validator.Validate(tcase.value)
		if len(tcase.expErrContains) == 0 {
			if err != nil {
				t.Fatalf("%d: %s", i+1, err.Error())
			}
		} else {
			if err == nil {
				t.Fatalf("%d: expected error got none", i+1)
			}
			for _, str := range tcase.expErrContains {
				if !strings.Contains(err.Error(), str) {
					t.Fatalf("%d: got %s, want %s", i+1, err.Error(), str)
				}
			}
		}
	}
}
