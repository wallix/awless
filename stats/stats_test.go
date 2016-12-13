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
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/wallix/awless/api"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
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

	awsInfra := &api.AwsInfra{}

	awsInfra.Instances = []*ec2.Instance{
		&ec2.Instance{InstanceId: aws.String("inst_1"), SubnetId: aws.String("sub_1"), VpcId: aws.String("vpc_1"), InstanceType: aws.String("t2.micro"), ImageId: aws.String("ami-e98bd29a")},
		&ec2.Instance{InstanceId: aws.String("inst_2"), SubnetId: aws.String("sub_2"), VpcId: aws.String("vpc_1"), InstanceType: aws.String("t2.micro"), ImageId: aws.String("ami-9398d3e0")},
		&ec2.Instance{InstanceId: aws.String("inst_3"), SubnetId: aws.String("sub_3"), VpcId: aws.String("vpc_2"), InstanceType: aws.String("t2.small"), ImageId: aws.String("ami-e98bd29a")},
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

	infra, err := rdf.BuildAwsInfraGraph("eu-west-1", awsInfra)
	if err != nil {
		t.Fatal(err)
	}

	awsAccess := &api.AwsAccess{}

	awsAccess.Groups = []*iam.Group{
		&iam.Group{GroupId: aws.String("group_1"), GroupName: aws.String("ngroup_1")},
		&iam.Group{GroupId: aws.String("group_2"), GroupName: aws.String("ngroup_2")},
	}

	awsAccess.LocalPolicies = []*iam.Policy{
		&iam.Policy{PolicyId: aws.String("policy_1"), PolicyName: aws.String("npolicy_1")},
		&iam.Policy{PolicyId: aws.String("policy_2"), PolicyName: aws.String("npolicy_2")},
	}

	awsAccess.Roles = []*iam.Role{
		&iam.Role{RoleId: aws.String("role_1")},
	}

	awsAccess.Users = []*iam.User{
		&iam.User{UserId: aws.String("usr_1")},
		&iam.User{UserId: aws.String("usr_2")},
		&iam.User{UserId: aws.String("usr_3")},
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

	access, err := rdf.BuildAwsAccessGraph("eu-west-1", awsAccess)
	if err != nil {
		t.Fatal(err)
	}

	db.AddLog("log msg 1")
	db.AddLog("log msg 2")
	db.AddLog("log msg 1")
	db.AddLog("log msg 3")
	db.AddLog("log msg 1")

	id, _ := db.GetStringValue(AWLESS_ID_KEY)
	aId, _ := db.GetStringValue(AWLESS_AID_KEY)
	expected := Stats{
		Id:      id,
		AId:     aId,
		Version: config.Version,
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
			Region:                "eu-west-1",
			NbVpcs:                2,
			NbSubnets:             3,
			NbInstances:           3,
			MinSubnetsPerVpc:      1,
			MaxSubnetsPerVpc:      2,
			MinInstancesPerSubnet: 1,
			MaxInstancesPerSubnet: 1,
		},
		InstancesStats: []*InstancesStat{
			{Type: "InstanceType", Date: now, Name: "t2.micro", Hits: 2},
			{Type: "InstanceType", Date: now, Name: "t2.small", Hits: 1},
			{Type: "ImageId", Date: now, Name: "ami-e98bd29a", Hits: 2},
			{Type: "ImageId", Date: now, Name: "ami-9398d3e0", Hits: 1},
		},
		AccessMetrics: &AccessMetrics{
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
		Logs: []*Log{
			{Msg: "log msg 1", Hits: 3, Date: now},
			{Msg: "log msg 2", Hits: 1, Date: now},
			{Msg: "log msg 3", Hits: 1, Date: now},
		},
	}

	stats, _, err := BuildStats(db, infra, access, 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Timestamps", func(t *testing.T) {
		for i := range expected.Commands {
			if got, want := stats.Commands[i].Date, expected.Commands[i].Date; !got.Equal(want) {
				t.Fatalf("got %v want %v", got, want)
			}
		}
		if got, want := stats.AccessMetrics.Date, expected.AccessMetrics.Date; !SameDay(&got, &want) {
			t.Fatalf("got %v want %v", got, want)
		}
		if got, want := stats.InfraMetrics.Date, expected.InfraMetrics.Date; !SameDay(&got, &want) {
			t.Fatalf("got %v want %v", got, want)
		}
		for i := range expected.InstancesStats {
			if got, want := stats.InstancesStats[i].Date, expected.InstancesStats[i].Date; !SameDay(&got, &want) {
				t.Fatalf("got %v want %v", got, want)
			}
		}
		for i := range expected.Logs {
			if got, want := stats.Logs[i].Date, expected.Logs[i].Date; !SameDay(&got, &want) {
				t.Fatalf("got %v want %v", got, want)
			}
		}
	})

	t.Run("Ignoring timestamps", func(t *testing.T) {
		sort.Sort(ByCommand(stats.Commands))
		sort.Sort(ByCommand(expected.Commands))
		sort.Sort(ByTypeAndName(stats.InstancesStats))
		sort.Sort(ByTypeAndName(expected.InstancesStats))
		nullifyTime(stats)
		nullifyTime(&expected)
		if got, want := reflect.DeepEqual(stats, &expected), true; got != want {
			t.Fatalf("got\n%+v\nwant\n%+v\n", *stats, expected)
		}
	})

	t.Run("SendStats", func(t *testing.T) {
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

		if err = db.SendStats(ts.URL, privateKey.PublicKey, infra, access); err != nil {
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

func (p *DailyCommands) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *InstancesStat) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *InfraMetrics) String() string {
	return fmt.Sprintf("%+v", *p)
}

func (p *AccessMetrics) String() string {
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
	case *Stats:
		nullifyTime(ii.AccessMetrics)
		nullifyTime(ii.Commands)
		nullifyTime(ii.InfraMetrics)
		nullifyTime(ii.InstancesStats)
		nullifyTime(ii.Logs)
	case *AccessMetrics, *InfraMetrics, *DailyCommands, *InstancesStat, *Log:
		reflect.ValueOf(i).Elem().FieldByName("Date").Set(reflect.ValueOf(time.Time{}))
	case []*DailyCommands:
		for _, v := range ii {
			nullifyTime(v)
		}
	case []*InstancesStat:
		for _, v := range ii {
			nullifyTime(v)
		}
	case []*Log:
		for _, v := range ii {
			nullifyTime(v)
		}
	default:
		panic(fmt.Sprintf("%T is not a known type", i))
	}
}

type ByCommand []*DailyCommands

func (a ByCommand) Len() int      { return len(a) }
func (a ByCommand) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool {
	return a[i].Command < a[j].Command
}

type ByTypeAndName []*InstancesStat

func (a ByTypeAndName) Len() int      { return len(a) }
func (a ByTypeAndName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTypeAndName) Less(i, j int) bool {
	if a[i].Type == a[j].Type {
		return a[i].Name < a[j].Name
	}
	return a[i].Type < a[j].Type
}
