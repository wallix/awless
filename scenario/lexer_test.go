package scenario

import (
	"testing"

	"github.com/wallix/awless/scenario/driver"
)

func TestParseScenarioLines(t *testing.T) {
	raw := `CREATE VPC CIDR 10.0.0.1/16 REF vpc_1
CREATE SUBNET REFERENCES vpc_1 REF subnet_1
CREATE INSTANCE COUNT 1 BASE linux REFERENCES subnet_1`

	lex := &Lexer{}

	scen := lex.ParseScenario(raw)

	if got, want := len(scen.Lines), 3; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if got, want := scen.Lines[0].Action, driver.CREATE; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := scen.Lines[0].Resource, driver.VPC; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := scen.Lines[0].Params[driver.CIDR], "10.0.0.1/16"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := scen.Lines[0].Params[driver.REF], "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	if got, want := scen.Lines[1].Action, driver.CREATE; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := scen.Lines[1].Resource, driver.SUBNET; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := scen.Lines[1].Params[driver.REF], "subnet_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := scen.Lines[1].Params[driver.REFERENCES], "vpc_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}

	if got, want := scen.Lines[2].Action, driver.CREATE; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := scen.Lines[2].Resource, driver.INSTANCE; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := scen.Lines[2].Params[driver.COUNT], "1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := scen.Lines[2].Params[driver.BASE], "linux"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
	if got, want := scen.Lines[2].Params[driver.REFERENCES], "subnet_1"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}
