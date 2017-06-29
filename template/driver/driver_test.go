package driver_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/driver"
)

func TestMultiDriver(t *testing.T) {
	ab := func(driver.Context, map[string]interface{}) (interface{}, error) { return "ab", nil }
	bc := func(driver.Context, map[string]interface{}) (interface{}, error) { return "bc", nil }
	de := func(driver.Context, map[string]interface{}) (interface{}, error) { return "de", nil }
	ef := func(driver.Context, map[string]interface{}) (interface{}, error) { return "ef", nil }
	mock1 := &mockDriver{
		lookupFn: func(lookups ...string) (driverFn driver.DriverFn, err error) {
			if len(lookups) != 1 {
				return nil, errors.New("expect only 1 param")
			}
			switch lookups[0] {
			case "ab":
				return ab, nil
			case "bc":
				return bc, nil
			case "de":
				return de, nil
			default:
				return nil, driver.ErrDriverFnNotFound
			}
		},
	}
	mock2 := &mockDriver{
		lookupFn: func(lookups ...string) (driverFn driver.DriverFn, err error) {
			if len(lookups) != 1 {
				return nil, errors.New("expect only 1 param")
			}
			switch lookups[0] {
			case "de":
				return de, nil
			case "ef":
				return ef, nil
			default:
				return nil, driver.ErrDriverFnNotFound
			}
		},
	}

	d := driver.NewMultiDriver(mock1, mock2)
	d.SetDryRun(true)
	if got, want := mock1.dryRun, true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	if got, want := mock2.dryRun, true; got != want {
		t.Fatalf("got %t, want %t", got, want)
	}
	d.SetLogger(logger.DiscardLogger)
	if got, want := mock1.logger, logger.DiscardLogger; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := mock2.logger, logger.DiscardLogger; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}

	tcases := []struct {
		str      string
		exectErr string
		expectFn driver.DriverFn
	}{
		{str: "ab", exectErr: "", expectFn: ab},
		{str: "bc", exectErr: "", expectFn: bc},
		{str: "cd", exectErr: "function corresponding to '[cd]' not found in drivers"},
		{str: "de", exectErr: "2 functions corresponding to '[de]' found in drivers"},
		{str: "ef", exectErr: "", expectFn: ef},
	}

	for i, tcase := range tcases {
		fn, err := d.Lookup(tcase.str)
		if err == nil {
			if got, want := "", tcase.exectErr; got != want {
				t.Fatalf("case %d, err: got %s, want %s", i, got, want)
			}
			if got, want := reflect.ValueOf(fn).Pointer(), reflect.ValueOf(tcase.expectFn).Pointer(); got != want {
				t.Fatalf("case %d, got %v, want %v", i, got, want)
			}
		} else {
			if got, want := err.Error(), tcase.exectErr; got != want {
				t.Fatalf("case %d, err: got %s, want %s", i, got, want)
			}
		}
	}

}

type mockDriver struct {
	dryRun   bool
	logger   *logger.Logger
	lookupFn func(lookups ...string) (driverFn driver.DriverFn, err error)
}

func (d *mockDriver) SetDryRun(dry bool)         { d.dryRun = dry }
func (d *mockDriver) SetLogger(l *logger.Logger) { d.logger = l }

func (d *mockDriver) Lookup(lookups ...string) (driverFn driver.DriverFn, err error) {
	return d.lookupFn(lookups...)
}
