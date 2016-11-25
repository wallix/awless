package stats

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestBuildStats(t *testing.T) {
	f, e := ioutil.TempFile(".", "test.db")
	if e != nil {
		t.Fatal(e)
	}
	defer os.Remove(f.Name())

	db, err := OpenDB(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.AddHistoryCommandWithTime([]string{"awless sync"}, time.Now().Add(-24*time.Hour))
	db.AddHistoryCommandWithTime([]string{"awless diff"}, time.Now().Add(-24*time.Hour))
	db.AddHistoryCommandWithTime([]string{"awless diff"}, time.Now().Add(-24*time.Hour))
	db.AddHistoryCommand([]string{"awless diff"})
	db.AddHistoryCommand([]string{"awless diff"})
	db.AddHistoryCommand([]string{"awless sync"})
	db.AddHistoryCommand([]string{"awless sync"})
	db.AddHistoryCommand([]string{"awless list instances"})
	db.AddHistoryCommand([]string{"awless list vpcs"})
	db.AddHistoryCommand([]string{"awless list subnets"})
	db.AddHistoryCommand([]string{"awless list instances"})

	stats, _, err := db.BuildStats(0)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(stats.DailyStats), 2; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
	yesterdayDate := time.Now().Add(-24 * time.Hour)
	if got, want := SameDay(&stats.DailyStats[0].Date, &yesterdayDate), true; got != want {
		t.Fatalf("got %t; want %t", got, want)
	}
	nowDate := time.Now()
	if got, want := SameDay(&stats.DailyStats[1].Date, &nowDate), true; got != want {
		t.Fatalf("got %t; want %t", got, want)
	}
	expectedYesterday :=
		map[string]int{
			"awless sync": 1,
			"awless diff": 2,
		}
	expectedToday :=
		map[string]int{
			"awless sync":           2,
			"awless diff":           2,
			"awless list instances": 2,
			"awless list vpcs":      1,
			"awless list subnets":   1,
		}

	if got, want := reflect.DeepEqual(expectedYesterday, stats.DailyStats[0].Commands), true; got != want {
		t.Fatalf("got \n%#v\n; want \n%#v", stats.DailyStats[0].Commands, expectedYesterday)
	}
	if got, want := reflect.DeepEqual(expectedToday, stats.DailyStats[1].Commands), true; got != want {
		t.Fatalf("got \n%#v\n; want \n%#v", stats.DailyStats[1].Commands, expectedToday)
	}
}

func TestSendStats(t *testing.T) {
	f, e := ioutil.TempFile(".", "test.db")
	if e != nil {
		t.Fatal(e)
	}
	defer os.Remove(f.Name())

	db, err := OpenDB(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	db.AddHistoryCommand([]string{"awless sync"})
	db.AddHistoryCommand([]string{"awless diff"})
	db.AddHistoryCommand([]string{"awless diff"})
	db.AddHistoryCommand([]string{"awless diff"})
	db.AddHistoryCommand([]string{"awless sync"})
	db.AddHistoryCommand([]string{"awless list instances"})
	db.AddHistoryCommand([]string{"awless list vpcs"})
	db.AddHistoryCommand([]string{"awless list subnets"})
	db.AddHistoryCommand([]string{"awless list instances"})

	expected, _, err := db.BuildStats(0)
	if err != nil {
		t.Fatal(err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	processed := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var encrypted EncryptedData
		if e := json.NewDecoder(r.Body).Decode(&encrypted); e != nil {
			t.Fatal(e)
			return
		}
		defer r.Body.Close()

		sessionKey, e := rsa.DecryptOAEP(sha256.New(), nil, privateKey, encrypted.Key, nil)
		if e != nil {
			t.Fatal(e)
			return
		}

		decrypted, e := aesDecrypt(encrypted.Data, sessionKey)
		if e != nil {
			t.Fatal(e)
			return
		}

		var received Stats
		if e := json.Unmarshal(decrypted, &received); e != nil {
			t.Fatal(e)
			return
		}

		assertEqual(t, &received, expected)
		processed = true

	}))
	defer ts.Close()

	if err := db.SendStats(ts.URL, privateKey.PublicKey); err != nil {
		t.Fatal(err)
	}

	if got, want := processed, true; got != want {
		t.Fatalf("got %t; want %t", got, want)
	}
}

func TestIfDataToSend(t *testing.T) {
	f, e := ioutil.TempFile(".", "test.db")
	if e != nil {
		t.Fatal(e)
	}
	defer os.Remove(f.Name())

	db, err := OpenDB(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if got, want := db.CheckStatsToSend(1*time.Hour), true; got != want {
		t.Fatalf("got %t; want %t", got, want)
	}

	db.SetTimeValue(SENT_TIME_KEY, time.Now().Add(-2*time.Hour))
	if got, want := db.CheckStatsToSend(1*time.Hour), true; got != want {
		t.Fatalf("got %t; want %t", got, want)
	}
	db.SetTimeValue(SENT_TIME_KEY, time.Now())

	if got, want := db.CheckStatsToSend(1*time.Hour), false; got != want {
		t.Fatalf("got %t; want %t", got, want)
	}

}

func assertEqual(t *testing.T, stats1, stats2 *Stats) {
	b1, err := json.Marshal(stats1)
	if err != nil {
		t.Fatal(err)
	}
	b2, err := json.Marshal(stats2)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(b1), string(b2); got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func aesDecrypt(encrypted, key []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	copy(nonce, encrypted)

	decrypted, err := gcm.Open(nil, nonce, encrypted[gcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
