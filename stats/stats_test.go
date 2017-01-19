package stats

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
)

func TestStats(t *testing.T) {
	db, close := newTestDb()
	defer close()

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

	awsInfra := &aws.AwsInfra{}

	awsInfra.Instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: awssdk.String("inst_1"), SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1"), InstanceType: awssdk.String("t2.micro"), ImageId: awssdk.String("ami-e98bd29a")},
		&ec2.Instance{InstanceId: awssdk.String("inst_2"), SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1"), InstanceType: awssdk.String("t2.micro"), ImageId: awssdk.String("ami-9398d3e0")},
		&ec2.Instance{InstanceId: awssdk.String("inst_3"), SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2"), InstanceType: awssdk.String("t2.small"), ImageId: awssdk.String("ami-e98bd29a")},
	}

	awsInfra.Vpcs = []*ec2.Vpc{
		&ec2.Vpc{VpcId: awssdk.String("vpc_1")},
		&ec2.Vpc{VpcId: awssdk.String("vpc_2")},
	}

	awsInfra.Subnets = []*ec2.Subnet{
		&ec2.Subnet{SubnetId: awssdk.String("sub_1"), VpcId: awssdk.String("vpc_1")},
		&ec2.Subnet{SubnetId: awssdk.String("sub_2"), VpcId: awssdk.String("vpc_1")},
		&ec2.Subnet{SubnetId: awssdk.String("sub_3"), VpcId: awssdk.String("vpc_2")},
	}

	infra, err := aws.BuildAwsInfraGraph("eu-west-1", awsInfra)
	if err != nil {
		t.Fatal(err)
	}

	awsAccess := &aws.AwsAccess{}

	awsAccess.Groups = []*iam.Group{
		&iam.Group{GroupId: awssdk.String("group_1"), GroupName: awssdk.String("ngroup_1")},
		&iam.Group{GroupId: awssdk.String("group_2"), GroupName: awssdk.String("ngroup_2")},
	}

	awsAccess.LocalPolicies = []*iam.Policy{
		&iam.Policy{PolicyId: awssdk.String("policy_1"), PolicyName: awssdk.String("npolicy_1")},
		&iam.Policy{PolicyId: awssdk.String("policy_2"), PolicyName: awssdk.String("npolicy_2")},
	}

	awsAccess.Roles = []*iam.Role{
		&iam.Role{RoleId: awssdk.String("role_1")},
	}

	awsAccess.Users = []*iam.User{
		&iam.User{UserId: awssdk.String("usr_1")},
		&iam.User{UserId: awssdk.String("usr_2")},
		&iam.User{UserId: awssdk.String("usr_3")},
	}

	awsAccess.UsersByGroup = map[string][]string{
		"group_1": []string{"usr_1", "usr_2"},
		"group_2": []string{"usr_1", "usr_2", "usr_3"},
	}

	awsAccess.UsersByLocalPolicies = map[string][]string{
		"policy_1": []string{"usr_1", "usr_2", "usr_3"},
		"policy_2": []string{"usr_1"},
	}

	awsAccess.RolesByLocalPolicies = map[string][]string{
		"policy_1": []string{"role_1"},
		"policy_2": []string{},
	}

	awsAccess.GroupsByLocalPolicies = map[string][]string{
		"policy_1": []string{"group_1"},
		"policy_2": []string{"group_1", "group_2"},
	}

	access, err := aws.BuildAwsAccessGraph("eu-west-1", awsAccess)
	if err != nil {
		t.Fatal(err)
	}

	db.AddLog("log msg 1")
	db.AddLog("log msg 2")
	db.AddLog("log msg 1")
	db.AddLog("log msg 3")
	db.AddLog("log msg 1")

	id, _ := db.GetStringValue(database.AwlessIdKey)
	aId, _ := db.GetStringValue(database.AwlessAIdKey)
	expected := stats{
		Id:      id,
		AId:     aId,
		Version: config.Version,
		Commands: []*dailyCommands{
			{Command: "awless sync", Hits: 1, Date: yesterday},
			{Command: "awless diff", Hits: 2, Date: yesterday},
			{Command: "awless diff", Hits: 1, Date: now},
			{Command: "awless sync", Hits: 1, Date: now},
			{Command: "awless list instances", Hits: 2, Date: now},
			{Command: "awless list vpcs", Hits: 1, Date: now},
		},
		InfraMetrics: &infraMetrics{
			Date:                  now,
			Region:                "eu-west-1",
			NbVpcs:                2,
			NbSubnets:             3,
			NbInstances:           3,
			MinSubnetsPerVpc:      1,
			MaxSubnetsPerVpc:      2,
			MinInstancesPerSubnet: 1,
			MaxInstancesPerSubnet: 1,
		},
		InstancesStats: []*instancesStat{
			{Type: "InstanceType", Date: now, Name: "t2.micro", Hits: 2},
			{Type: "InstanceType", Date: now, Name: "t2.small", Hits: 1},
			{Type: "ImageId", Date: now, Name: "ami-e98bd29a", Hits: 2},
			{Type: "ImageId", Date: now, Name: "ami-9398d3e0", Hits: 1},
		},
		AccessMetrics: &accessMetrics{
			Date:                     now,
			Region:                   "eu-west-1",
			NbGroups:                 2,
			NbPolicies:               2,
			NbRoles:                  1,
			NbUsers:                  3,
			MinUsersByGroup:          2,
			MaxUsersByGroup:          3,
			MinUsersByLocalPolicies:  1,
			MaxUsersByLocalPolicies:  3,
			MinRolesByLocalPolicies:  0,
			MaxRolesByLocalPolicies:  1,
			MinGroupsByLocalPolicies: 1,
			MaxGroupsByLocalPolicies: 2,
		},
		Logs: []*database.Log{
			{Msg: "log msg 1", Hits: 3, Date: now},
			{Msg: "log msg 2", Hits: 1, Date: now},
			{Msg: "log msg 3", Hits: 1, Date: now},
		},
	}

	s, _, err := BuildStats(db, infra, access, 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Timestamps", func(t *testing.T) {
		if got, want := len(s.Commands), len(expected.Commands); got != want {
			t.Fatalf("got %d want %d", got, want)
		}
		for i := range expected.Commands {
			if got, want := s.Commands[i].Date, expected.Commands[i].Date; !got.Equal(want) {
				t.Fatalf("got %v want %v", got, want)
			}
		}
		if got, want := s.AccessMetrics.Date, expected.AccessMetrics.Date; !sameDay(&got, &want) {
			t.Fatalf("got %v want %v", got, want)
		}
		if got, want := s.InfraMetrics.Date, expected.InfraMetrics.Date; !sameDay(&got, &want) {
			t.Fatalf("got %v want %v", got, want)
		}
		if got, want := len(s.InstancesStats), len(expected.InstancesStats); got != want {
			t.Fatalf("got %d want %d", got, want)
		}
		for i := range expected.InstancesStats {
			if got, want := s.InstancesStats[i].Date, expected.InstancesStats[i].Date; !sameDay(&got, &want) {
				t.Fatalf("got %v want %v", got, want)
			}
		}
		for i := range expected.Logs {
			if got, want := s.Logs[i].Date, expected.Logs[i].Date; !sameDay(&got, &want) {
				t.Fatalf("got %v want %v", got, want)
			}
		}
	})

	t.Run("Ignoring timestamps", func(t *testing.T) {
		sort.Sort(ByCommand(s.Commands))
		sort.Sort(ByCommand(expected.Commands))
		sort.Sort(ByTypeAndName(s.InstancesStats))
		sort.Sort(ByTypeAndName(expected.InstancesStats))
		nullifyTime(s)
		nullifyTime(&expected)
		if got, want := reflect.DeepEqual(s, &expected), true; got != want {
			t.Fatalf("got\n%+v\nwant\n%+v\n", *s, expected)
		}
	})

	t.Run("SendStats", func(t *testing.T) {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}
		publicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			t.Fatal(err)
		}
		serverPublicKey = string(pem.EncodeToMemory(
			&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: publicKey,
			},
		))

		processed := false

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var encrypted encryptedData
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

			var received stats
			if e := json.NewDecoder(extracted).Decode(&received); e != nil {
				t.Fatal(e)
				return
			}

			sort.Sort(ByCommand(received.Commands))
			sort.Sort(ByCommand(expected.Commands))
			sort.Sort(ByTypeAndName(received.InstancesStats))
			sort.Sort(ByTypeAndName(expected.InstancesStats))
			nullifyTime(&received)
			nullifyTime(&expected)

			if !reflect.DeepEqual(received, expected) {
				t.Fatalf("got %+v; want %+v", received, expected)
			}
			processed = true

		}))
		defer ts.Close()
		serverUrl = ts.URL

		if err = SendStats(db, infra, access); err != nil {
			t.Fatal(err)
		}

		if got, want := processed, true; got != want {
			t.Fatalf("got %t; want %t", got, want)
		}

		logs, err := db.GetLogs()
		if err != nil {
			t.Fatal(err)
		}

		if got, want := len(logs), 0; got != want {
			t.Fatalf("got %d; want %d", got, want)
		}
	})
}

func TestIfDataToSend(t *testing.T) {
	db, close := newTestDb()
	defer close()

	expirationDuration = 1 * time.Hour
	if got, want := CheckStatsToSend(db), true; got != want {
		t.Fatalf("got %t; want %t", got, want)
	}

	db.SetTimeValue(database.SentTimeKey, time.Now().Add(-2*time.Hour))
	if got, want := CheckStatsToSend(db), true; got != want {
		t.Fatalf("got %t; want %t", got, want)
	}
	db.SetTimeValue(database.SentTimeKey, time.Now())

	if got, want := CheckStatsToSend(db), false; got != want {
		t.Fatalf("got %t; want %t", got, want)
	}
}

func (p *dailyCommands) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *instancesStat) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *infraMetrics) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *accessMetrics) String() string {
	return fmt.Sprintf("%+v", *p)
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

func nullifyTime(i interface{}) {
	switch ii := i.(type) {
	case *stats:
		nullifyTime(ii.AccessMetrics)
		nullifyTime(ii.Commands)
		nullifyTime(ii.InfraMetrics)
		nullifyTime(ii.InstancesStats)
		nullifyTime(ii.Logs)
	case *accessMetrics, *infraMetrics, *dailyCommands, *instancesStat, *database.Log:
		reflect.ValueOf(i).Elem().FieldByName("Date").Set(reflect.ValueOf(time.Time{}))
	case []*dailyCommands:
		for _, v := range ii {
			nullifyTime(v)
		}
	case []*instancesStat:
		for _, v := range ii {
			nullifyTime(v)
		}
	case []*database.Log:
		for _, v := range ii {
			nullifyTime(v)
		}
	default:
		panic(fmt.Sprintf("%T is not a known type", i))
	}
}

type ByCommand []*dailyCommands

func (a ByCommand) Len() int      { return len(a) }
func (a ByCommand) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool {
	return a[i].Command < a[j].Command
}

type ByTypeAndName []*instancesStat

func (a ByTypeAndName) Len() int      { return len(a) }
func (a ByTypeAndName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTypeAndName) Less(i, j int) bool {
	if a[i].Type == a[j].Type {
		return a[i].Name < a[j].Name
	}
	return a[i].Type < a[j].Type
}
