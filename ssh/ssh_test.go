package ssh

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	gossh "golang.org/x/crypto/ssh"
)

func TestCheckHostKey(t *testing.T) {
	//Create env
	f, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(f)
	err = os.Setenv("HOME", f)
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(filepath.Join(f, ".ssh"), 0755)
	if err != nil {
		t.Fatal(err)
	}
	knownHostsFileContent := `1.2.3.4 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBFJz/HFJUq6SaXD5FdLe6ddIpmNPFim7E3NkNCSNurDun/h3BOIzNGfuseyMn32n24oQayhjkX8eGqevJIA38E=
3.4.5.6 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBKl6fXNb/yA0w7brzqNuOCwLJ/aPEMerl7/lsF0Y/1oafD2bxzj+QsEZo4XK/kvwCjqQArFO5nET+Tz015C6Kk=
4.5.6.7 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBKl6fXNb/yA0w7brzqNuOCwLJ/aPEMerl7/lsF0Y/1oafD2bxzj+QsEZo4XK/kvwCjqQArFO5nET+Tz015C6Kk=
`
	knowHostsFile := filepath.Join(f, ".ssh", "known_hosts")
	err = ioutil.WriteFile(knowHostsFile, []byte(knownHostsFileContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	ecdsa1, _, _, _, err := gossh.ParseAuthorizedKey([]byte("1.2.3.4 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBFJz/HFJUq6SaXD5FdLe6ddIpmNPFim7E3NkNCSNurDun/h3BOIzNGfuseyMn32n24oQayhjkX8eGqevJIA38E="))
	if err != nil {
		t.Fatalf("error parsing ecdsa1: %s", err)
	}
	ecdsa2, _, _, _, err := gossh.ParseAuthorizedKey([]byte("3.4.5.6 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBKl6fXNb/yA0w7brzqNuOCwLJ/aPEMerl7/lsF0Y/1oafD2bxzj+QsEZo4XK/kvwCjqQArFO5nET+Tz015C6Kk="))
	if err != nil {
		t.Fatalf("error parsing ecdsa1: %s", err)
	}
	knownKeys := make(map[string]gossh.PublicKey)
	numberKeysAdded := 0
	trustKeyFunc = func(hostname string, remote net.Addr, key gossh.PublicKey, _ string) bool {
		knownKeys[hostname] = key
		numberKeysAdded++
		return true
	}
	tcases := []struct {
		ip     string
		key    gossh.PublicKey
		expErr string
	}{
		{"1.2.3.4", ecdsa1, ""},
		{"2.3.4.5", ecdsa1, ""},
		{"2.3.4.5", ecdsa1, ""},
		{"3.4.5.6", ecdsa1, fmt.Sprintf(`
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
AWLESS DETECTED THAT THE REMOTE HOST PUBLIC KEY HAS CHANGED
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

Host key for '3.4.5.6:22' has changed and you did not disable strict host key checking.
Someone may be trying to intercept your connection (man-in-the-middle attack). Otherwise, the host key may have been changed.

The fingerprint for the ecdsa-sha2-nistp256 key sent by the remote host is SHA256:t88p7xhU/1D3USAczv1d88hTZJbOeWN/ktcNmeWh6qI.
You persisted:
-> SHA256:rEoUFSq4TQYk6NO00isOyR7s/4keg0NePYE39XIf5vY (ecdsa-sha2-nistp256 key in %s:2)

To get rid of this message, update '%[1]s:2'`, knowHostsFile)},
		{"3.4.5.6", ecdsa2, ""},
	}

	for i, tcase := range tcases {
		addr, er := net.ResolveTCPAddr("", tcase.ip+":22")
		if er != nil {
			t.Fatal(er)
		}
		got := checkHostKey(tcase.ip+":22", addr, tcase.key)
		var gotStr string
		if got != nil {
			gotStr = got.Error()
		}
		if got, want := gotStr, tcase.expErr; got != want {
			t.Fatalf("case %d: got '%s', want '%s'", i+1, got, want)
		}
	}
	expectedKnownKeys := map[string]gossh.PublicKey{
		"2.3.4.5:22": ecdsa1,
	}
	if got, want := numberKeysAdded, 1; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := knownKeys, expectedKnownKeys; !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}

	knownHostContent, err := ioutil.ReadFile(knowHostsFile)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(knownHostContent), knownHostsFileContent+"2.3.4.5 ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBFJz/HFJUq6SaXD5FdLe6ddIpmNPFim7E3NkNCSNurDun/h3BOIzNGfuseyMn32n24oQayhjkX8eGqevJIA38E=\n"; got != want {
		t.Fatalf("got \n%s\nwant\n%s\n", got, want)
	}
}
