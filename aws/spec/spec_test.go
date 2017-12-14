package awsspec

import (
	"strings"
	"testing"

	"github.com/wallix/awless/template/env"
	"github.com/wallix/awless/template/params"
)

func checkErrs(t *testing.T, errs []error, length int, expected ...string) {
	t.Helper()
	if got, want := len(errs), length; got != want {
		t.Fatalf("got %d want %d", got, want)
	}
	for _, exp := range expected {
		var found bool
		for _, err := range errs {
			if strings.Contains(err.Error(), exp) {
				found = true
			}
		}
		if !found {
			t.Fatalf("'%s' should be in errors %v", exp, errs)
		}
	}
}

type validParamTestStruct struct {
	Param1 *string `templateName:"param1" required:""`
	Param2 *string `templateName:"param2" required:""`
	Param3 *string `templateName:"param3"`
}

func (*validParamTestStruct) Params() params.Rule                                      { return nil }
func (*validParamTestStruct) ValidateCommand(map[string]interface{}, []string) []error { return nil }
func (*validParamTestStruct) inject(params map[string]interface{}) error               { return nil }
func (*validParamTestStruct) Run(renv env.Running, params map[string]interface{}) (interface{}, error) {
	return nil, nil
}

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

func keys(m map[string]interface{}) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return
}

func convertParamsIfAvailable(cmd interface{}, params map[string]interface{}) (map[string]interface{}, error) {
	type C interface {
		ConvertParams() ([]string, func(values map[string]interface{}) (map[string]interface{}, error))
	}
	if v, ok := cmd.(C); ok {
		keys, convFunc := v.ConvertParams()
		values := make(map[string]interface{})
		for _, k := range keys {
			if vv, ok := params[k]; ok {
				values[k] = vv
			}
		}
		converted, err := convFunc(values)
		if err != nil {
			return params, err
		}
		for _, k := range keys {
			delete(params, k)
		}
		for k, v := range converted {
			params[k] = v
		}
	}
	return params, nil
}
