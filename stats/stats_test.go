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
	"github.com/aws/aws-sdk-go/service/iam"
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

	infra, err := rdf.BuildAwsInfraGraph("eu-west-1", awsInfra)
	if err != nil {
		t.Fatal(err)
	}

	id, _ := db.GetStringValue(AWLESS_ID_KEY)
	expected := Stats{
		Id:      id,
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
			Region:                "",
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
			Region:                   "",
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
	}

	access, err := rdf.BuildAwsAccessGraph("eu-west-1", awsAccess)
	if err != nil {
		t.Fatal(err)
	}

	stats, _, err := BuildStats(db, infra, access, 0)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(stats.Commands), len(expected.Commands); got != want {
		t.Fatalf("got %d; want %d", got, want)
	}

	if got, want := statsEqual(stats, &expected), true; got != want {
		t.Fatalf("got\n%+v\nwant\n%+v\n", *stats, expected)
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

	localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.InfraFilename))
	if err != nil {
		t.Fatal(err)
	}

	localAccess, err := rdf.NewGraphFromFile(filepath.Join(config.GitDir, config.AccessFilename))
	if err != nil {
		t.Fatal(err)
	}

	expected, _, err := BuildStats(db, localInfra, localAccess, 0)
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
	if stats1.Version != stats2.Version {
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
	return infraMetricsEqual(stats1.InfraMetrics, stats2.InfraMetrics) &&
		instancesStatsEqual(stats1.InstancesStats, stats2.InstancesStats) &&
		accessMetricsEqual(stats1.AccessMetrics, stats2.AccessMetrics)
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

func accessMetricsEqual(a1, a2 *AccessMetrics) bool {
	if a1 == a2 {
		return true
	}
	if a1 == nil {
		return false
	}

	if a1.Region != a2.Region {
		return false
	}
	if a1.NbGroups != a2.NbGroups || a1.NbPolicies != a2.NbPolicies || a1.NbRoles != a2.NbRoles || a1.NbUsers != a2.NbUsers {
		return false
	}
	if a1.MinUsersByGroup != a2.MinUsersByGroup || a1.MaxUsersByGroup != a2.MaxUsersByGroup {
		return false
	}
	if a1.MinUsersByLocalPolicies != a2.MinUsersByLocalPolicies || a1.MaxUsersByLocalPolicies != a2.MaxUsersByLocalPolicies {
		return false
	}
	if a1.MinRolesByLocalPolicies != a2.MinRolesByLocalPolicies || a1.MaxRolesByLocalPolicies != a2.MaxRolesByLocalPolicies {
		return false
	}

	if a1.MinGroupsByLocalPolicies != a2.MinGroupsByLocalPolicies || a1.MaxGroupsByLocalPolicies != a2.MaxGroupsByLocalPolicies {
		return false
	}
	return SameDay(&a1.Date, &a2.Date)
}

func instancesStatsEqual(is1, is2 []*InstancesStat) bool {
	if is1 == nil && is2 == nil {
		return true
	}
	if is1 == nil {
		return false
	}
	if len(is1) != len(is2) {
		return false
	}
	for _, i1 := range is1 {
		found := false
		for _, i2 := range is2 {
			if instancesStatEqual(i1, i2) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func instancesStatEqual(i1, i2 *InstancesStat) bool {
	if i1 == i2 {
		return true
	}
	if i1 == nil || i2 == nil {
		return false
	}
	return SameDay(&i1.Date, &i2.Date) && i1.Hits == i2.Hits && i1.Name == i2.Name && i1.Type == i2.Type
}

func assertEqual(t *testing.T, got, want *Stats) {
	if !statsEqual(got, want) {
		t.Fatalf("got %+v; want %+v", got, want)
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
