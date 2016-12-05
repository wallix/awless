package stats

import (
	"bytes"
	"compress/gzip"
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
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/wallix/awless/api"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
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
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	db.AddHistoryCommandWithTime([]string{"awless sync"}, yesterday)
	db.AddHistoryCommandWithTime([]string{"awless diff"}, yesterday)
	db.AddHistoryCommandWithTime([]string{"awless diff"}, yesterday)
	db.AddHistoryCommandWithTime([]string{"awless diff"}, now)
	db.AddHistoryCommandWithTime([]string{"awless sync"}, now)
	db.AddHistoryCommandWithTime([]string{"awless list instances"}, now)
	db.AddHistoryCommandWithTime([]string{"awless list vpcs"}, now)
	db.AddHistoryCommandWithTime([]string{"awless list instances"}, now)

	awsInfra := &api.AwsInfra{}

	awsInfra.Instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: aws.String("inst_1"), SubnetId: aws.String("sub_1"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_2"), SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1")},
		&ec2.Instance{InstanceId: aws.String("inst_3"), SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_2")},
	}

	awsInfra.Vpcs = []*ec2.Vpc{
		&ec2.Vpc{VpcId: aws.String("vpc_1")},
		&ec2.Vpc{VpcId: aws.String("vpc_2")},
	}

	awsInfra.Subnets = []*ec2.Subnet{
		&ec2.Subnet{SubnetId: aws.String("sub_1"), VpcId: aws.String("vpc_1")},
		&ec2.Subnet{SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1")},
		&ec2.Subnet{SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_2")},
	}

	infra, _, err := rdf.BuildAwsInfraGraph("eu-west-1", awsInfra)
	if err != nil {
		t.Fatal(err)
	}

	id, _ := db.GetStringValue(AWLESS_ID_KEY)
	expected := Stats{
		Id: id,
		Commands: []*DailyCommands{
			{Command: "awless sync", Hits: 1, Date: yesterday},
			{Command: "awless diff", Hits: 2, Date: yesterday},
			{Command: "awless diff", Hits: 1, Date: now},
			{Command: "awless sync", Hits: 1, Date: now},
			{Command: "awless list instances", Hits: 2, Date: now},
			{Command: "awless list vpcs", Hits: 1, Date: now},
		},
		InfraMetrics: &InfraMetrics{
			Date:                  now,
			Region:                "",
			NbVpcs:                2,
			NbSubnets:             3,
			NbInstances:           3,
			MinSubnetsPerVpc:      1,
			MaxSubnetsPerVpc:      2,
			MinInstancesPerSubnet: 1,
			MaxInstancesPerSubnet: 1,
		},
	}

	stats, _, err := BuildStats(db, infra, 0)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(stats.Commands), len(expected.Commands); got != want {
		t.Fatalf("got %d; want %d", got, want)
	}

	if got, want := statsEqual(stats, &expected), true; got != want {
		t.Fatalf("got %#v; want %#v", *stats, expected)
	}
}

func TestBuildMetrics(t *testing.T) {
	infra, err := rdf.NewGraphFromFile("testdata/infra.rdf")
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	infraMetrics, err := buildInfraMetrics(infra, now)
	if err != nil {
		t.Fatal(err)
	}

	expectedMetrics := &InfraMetrics{
		Date:                  now,
		Region:                "",
		NbVpcs:                3,
		NbSubnets:             7,
		MinSubnetsPerVpc:      2,
		MaxSubnetsPerVpc:      3,
		NbInstances:           18,
		MinInstancesPerSubnet: 0,
		MaxInstancesPerSubnet: 4,
	}

	if got, want := reflect.DeepEqual(infraMetrics, expectedMetrics), true; got != want {
		t.Fatalf("got \n%#v\n; want \n%#v\n", infraMetrics, expectedMetrics)
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

	localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.Dir, config.InfraFilename))
	if err != nil {
		t.Fatal(err)
	}

	expected, _, err := BuildStats(db, localInfra, 0)
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

		extracted, e := gzip.NewReader(bytes.NewReader(decrypted))
		if e != nil {
			t.Fatal(e)
			return
		}
		defer extracted.Close()

		var received Stats
		if e := json.NewDecoder(extracted).Decode(&received); e != nil {
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

func dcEqual(dc1, dc2 *DailyCommands) bool {
	if dc1 == dc2 {
		return true
	}
	if dc1 == nil {
		return false
	}
	return dc1.Command == dc2.Command && dc1.Date.Equal(dc2.Date) && dc1.Hits == dc2.Hits
}

func statsEqual(stats1, stats2 *Stats) bool {
	if stats1 == stats2 {
		return true
	}
	if stats1 == nil {
		return false
	}
	if stats1.Id != stats2.Id {
		return false
	}
	for _, dc1 := range stats1.Commands {
		found := false
		for _, dc2 := range stats2.Commands {
			if dcEqual(dc1, dc2) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return infraMetricsEqual(stats1.InfraMetrics, stats2.InfraMetrics)
}

func infraMetricsEqual(i1, i2 *InfraMetrics) bool {
	if i1 == i2 {
		return true
	}
	if i1 == nil {
		return false
	}

	if i1.Region != i2.Region {
		return false
	}
	if i1.NbVpcs != i2.NbVpcs || i1.MaxSubnetsPerVpc != i2.MaxSubnetsPerVpc || i1.MinSubnetsPerVpc != i2.MinSubnetsPerVpc {
		return false
	}
	if i1.NbSubnets != i2.NbSubnets || i1.MaxInstancesPerSubnet != i2.MaxInstancesPerSubnet || i1.MinInstancesPerSubnet != i2.MinInstancesPerSubnet {
		return false
	}
	if i1.NbInstances != i2.NbInstances {
		return false
	}
	return SameDay(&i1.Date, &i2.Date)
}

func assertEqual(t *testing.T, stats1, stats2 *Stats) {
	if got, want := statsEqual(stats1, stats2), true; got != want {
		t.Fatalf("got %t; want %t", got, want)
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
