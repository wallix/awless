package aws

import "testing"

func TestParseImageQueryString(t *testing.T) {
	tcases := []struct {
		in  string
		out imageQuery
	}{
		{
			in:  "Canonical:ubuntu",
			out: imageQuery{Platform: Platforms["canonical"], Distro: Distro{Name: "ubuntu", Variant: Platforms["canonical"].LatestVariant, Arch: defaultArch, Virt: defaultVirt, Store: defaultStore}}},
		{
			in:  "canonical:ubuntu[jessie]",
			out: imageQuery{Platform: Platforms["canonical"], Distro: Distro{Name: "ubuntu", Variant: "jessie", Arch: defaultArch, Virt: defaultVirt, Store: defaultStore}}},
		{
			in:  "canonical:ubuntu[jessie]:i386:partial:instance-store",
			out: imageQuery{Platform: Platforms["canonical"], Distro: Distro{Name: "ubuntu", Variant: "jessie", Arch: "i386", Virt: "partial", Store: "instance-store"}}},
		{
			in:  "redhat:RHEL[7.3]:::instance-store",
			out: imageQuery{Platform: Platforms["redhat"], Distro: Distro{Name: "rhel", Variant: "7.3", Arch: defaultArch, Virt: defaultVirt, Store: "instance-store"}}},
		{
			in:  "debian",
			out: imageQuery{Platform: Platforms["debian"], Distro: Distro{Name: "debian", Variant: "jessie", Arch: defaultArch, Virt: defaultVirt, Store: defaultStore}}},
	}

	for _, tcase := range tcases {
		q, err := parseQuery(tcase.in)
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
