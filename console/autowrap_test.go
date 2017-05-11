package console

import "testing"

func TestAutoWrap(t *testing.T) {
	tcases := []struct {
		input        string
		maxWith      int
		wrappingChar string
		output       string
	}{
		{"myverylonglinewithoutanyspace", 6, "", "myvery longli newith outany space"},
		{"my very long line with spaces", 6, "\n", "my \nvery \nlong \nline \nwith \nspaces"},
		{"my very long line with spaces", 8, "\n", "my very \nlong \nline \nwith \nspaces"},
		{"my:very:very:very:long:arn", 8, " ", "my:very: very: very: long:arn"},
		{"my:very:very:very:long:arn", 16, " ", "my:very:very: very:long:arn"},
		{"nosplit", 16, " ", "nosplit"},
		{"splitateachchar", 1, " ", "s p l i t a t e a c h c h a r"},
		{"four;char+word.sepa+rate:with/spec+ials!char", 5, " ", "four; char+ word. sepa+ rate: with/ spec+ ials! char"},
		{"four;char+word.sepa+rate:with/spec+ials!char", 7, " ", "four; char+ word. sepa+ rate: with/ spec+ ials! char"},
		{"test of a long string with lots of spaces", 7, " ", "test of a long string with lots of spaces"},
	}
	for _, tcase := range tcases {
		wraper := autoWraper{maxWidth: tcase.maxWith, wrappingChar: tcase.wrappingChar}
		if got, want := wraper.Wrap(tcase.input), tcase.output; got != want {
			t.Fatalf("got '%q', want '%q'", got, want)
		}
	}
}
