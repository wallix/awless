package awsspec

import (
	"strings"
	"testing"
)

func TestParseImageQueryErrorCases(t *testing.T) {
	tcases := []struct {
		in          string
		errContains []string
	}{
		{in: "canonical:ubuntu:xenial:x86_64:hvm:ebs:", errContains: []string{"tokens"}},
		{in: "canonical:::::wrong-store", errContains: []string{"ebs", "instance-store"}},
		{in: "canonical::::wrong-virt", errContains: []string{"hvm", "paravirtual"}},
		{in: "canonical:::wrong-arch", errContains: []string{"i386", "x86_64"}},
		{in: "wrong", errContains: []string{"supported", "canonical"}},
	}

	for _, tcase := range tcases {
		_, err := ParseImageQuery(tcase.in)
		if err == nil {
			t.Fatalf("parsing '%s': expecting error got none", tcase.in)
		}

		for _, s := range tcase.errContains {
			if msg := err.Error(); !strings.Contains(msg, s) {
				t.Errorf("expecting %s to contain %s", msg, s)
			}

		}
	}
}

func TestImageQueryToString(t *testing.T) {
	tcases := []struct {
		in  string
		out string
	}{
		{in: "canonical", out: "canonical:ubuntu:xenial:x86_64:hvm:ebs"},
		{in: "canonical:ubuntu:trusty::paravirtual:ebs", out: "canonical:ubuntu:trusty:x86_64:paravirtual:ebs"},
		{in: "canonical:ubuntu:trusty:i386:paravirtual:instance-store", out: "canonical:ubuntu:trusty:i386:paravirtual:instance-store"},
	}

	for _, tcase := range tcases {
		q, err := ParseImageQuery(tcase.in)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := q.String(), tcase.out; got != want {
			t.Fatalf("got %s, want %s", got, want)
		}
	}
}

func TestParseImageQueryString(t *testing.T) {
	tcases := []struct {
		in  string
		out ImageQuery
	}{
		{
			in:  "Canonical:ubuntu",
			out: ImageQuery{Platform: Platforms["canonical"], Distro: Distro{Name: "ubuntu", Variant: Platforms["canonical"].LatestVariant, Arch: defaultArch, Virt: defaultVirt, Store: defaultStore}}},
		{
			in:  "Canonical:ubuntu::::",
			out: ImageQuery{Platform: Platforms["canonical"], Distro: Distro{Name: "ubuntu", Variant: Platforms["canonical"].LatestVariant, Arch: defaultArch, Virt: defaultVirt, Store: defaultStore}}},
		{
			in:  "canonical:ubuntu:jessie",
			out: ImageQuery{Platform: Platforms["canonical"], Distro: Distro{Name: "ubuntu", Variant: "jessie", Arch: defaultArch, Virt: defaultVirt, Store: defaultStore}}},
		{
			in:  "canonical:ubuntu:jessie:i386:paravirtual:instance-store",
			out: ImageQuery{Platform: Platforms["canonical"], Distro: Distro{Name: "ubuntu", Variant: "jessie", Arch: "i386", Virt: "paravirtual", Store: "instance-store"}}},
		{
			in:  "redhat:RHEL:7.3:::instance-store",
			out: ImageQuery{Platform: Platforms["redhat"], Distro: Distro{Name: "rhel", Variant: "7.3", Arch: defaultArch, Virt: defaultVirt, Store: "instance-store"}}},
		{
			in:  "centos:centos:7",
			out: ImageQuery{Platform: Platforms["centos"], Distro: Distro{Name: "centos", Variant: "7", Arch: defaultArch, Virt: defaultVirt, Store: defaultStore}}},
		{
			in:  "debian",
			out: ImageQuery{Platform: Platforms["debian"], Distro: Distro{Name: "debian", Variant: "stretch", Arch: defaultArch, Virt: defaultVirt, Store: defaultStore}}},
	}

	for _, tcase := range tcases {
		q, err := ParseImageQuery(tcase.in)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := q.Platform.Id, tcase.out.Platform.Id; got != want {
			t.Fatalf("parsing %s: query platform id: got %v, want %v", tcase.in, got, want)
		}
		if got, want := q.Distro, tcase.out.Distro; got != want {
			t.Fatalf("parsing %s: query distro: got %v, want %v", tcase.in, got, want)
		}
	}
}
