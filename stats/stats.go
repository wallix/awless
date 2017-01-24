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
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/wallix/awless/config"
	"github.com/wallix/awless/database"
	"github.com/wallix/awless/graph"
)

var (
	serverUrl          = "http://52.213.243.16"
	expirationDuration = 24 * time.Hour
)

func SendStats(db *database.DB, localInfra, localAccess *graph.Graph) error {
	publicKey, err := loadPublicKey()
	if err != nil {
		return err
	}
	lastCommandId, err := db.GetIntValue(database.SentIdKey)
	if err != nil {
		return err
	}

	s, lastCommandId, err := BuildStats(db, localInfra, localAccess, lastCommandId)
	if err != nil {
		return err
	}

	var zipped bytes.Buffer
	zippedW := gzip.NewWriter(&zipped)
	if err = json.NewEncoder(zippedW).Encode(s); err != nil {
		return err
	}
	zippedW.Close()

	sessionKey, encrypted, err := aesEncrypt(zipped.Bytes())
	if err != nil {
		return err
	}
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, sessionKey, nil)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(encryptedData{encryptedKey, encrypted})
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 2 * time.Second}
	if _, err := client.Post(serverUrl, "application/json", bytes.NewReader(payload)); err != nil {
		return err
	}

	if err := db.SetIntValue(database.SentIdKey, lastCommandId); err != nil {
		return err
	}
	if err := db.SetTimeValue(database.SentTimeKey, time.Now()); err != nil {
		return err
	}
	if err := db.DeleteLogs(); err != nil {
		return err
	}
	if err := db.DeleteHistory(); err != nil {
		return err
	}
	return nil
}

func BuildStats(db *database.DB, infra *graph.Graph, access *graph.Graph, fromCommandId int) (*stats, int, error) {
	commandsStat, lastCommandId, err := buildCommandsStat(db, fromCommandId)
	if err != nil {
		return nil, 0, err
	}
	region := db.MustGetDefaultRegion()

	im := &infraMetrics{}
	if infra != nil {
		im, err = buildInfraMetrics(region, infra)
		if err != nil {
			return nil, 0, err
		}
	}

	am := &accessMetrics{}
	if access != nil {
		am, err = buildAccessMetrics(region, access, time.Now())
		if err != nil {
			return nil, 0, err
		}
	}

	is, err := buildInstancesStats(infra)
	if err != nil {
		return nil, 0, err
	}

	id, err := db.GetStringValue(database.AwlessIdKey)
	if err != nil {
		return nil, 0, err
	}

	aId, err := db.GetStringValue(database.AwlessAIdKey)
	if err != nil {
		return nil, 0, err
	}

	logs, err := db.GetLogs()
	if err != nil {
		return nil, 0, err
	}

	s := &stats{
		Id:             id,
		AId:            aId,
		Version:        config.Version,
		BuildInfo:      config.CurrentBuildInfo,
		Commands:       commandsStat,
		InfraMetrics:   im,
		InstancesStats: is,
		AccessMetrics:  am,
		Logs:           logs,
	}

	return s, lastCommandId, nil
}

func CheckStatsToSend(db *database.DB) bool {
	sent, err := db.GetTimeValue(database.SentTimeKey)
	if err != nil {
		sent = time.Time{}
	}
	return (time.Since(sent) > expirationDuration)
}

type stats struct {
	Id             string
	AId            string
	Version        string
	BuildInfo      config.BuildInfo
	Commands       []*dailyCommands
	InfraMetrics   *infraMetrics
	InstancesStats []*instancesStat
	AccessMetrics  *accessMetrics
	Logs           []*database.Log
}

type dailyCommands struct {
	Command string
	Hits    int
	Date    time.Time
}

type instancesStat struct {
	Type string
	Date time.Time
	Hits int
	Name string
}

type accessMetrics struct {
	Date                     time.Time
	Region                   string
	NbGroups                 int
	NbPolicies               int
	NbRoles                  int
	NbUsers                  int
	MinUsersByGroup          int
	MaxUsersByGroup          int
	MinUsersByLocalPolicies  int
	MaxUsersByLocalPolicies  int
	MinRolesByLocalPolicies  int
	MaxRolesByLocalPolicies  int
	MinGroupsByLocalPolicies int
	MaxGroupsByLocalPolicies int
}

func buildCommandsStat(db *database.DB, fromCommandId int) ([]*dailyCommands, int, error) {
	var commandsStat []*dailyCommands

	commandsHistory, err := db.GetHistory(fromCommandId)
	if err != nil {
		return commandsStat, 0, err
	}

	if len(commandsHistory) == 0 {
		return commandsStat, 0, nil
	}

	date := commandsHistory[0].Time
	commands := make(map[string]int)

	for _, command := range commandsHistory {
		if !sameDay(&date, &command.Time) {
			commandsStat = addDailyCommands(commandsStat, commands, &date)
			date = command.Time
			commands = make(map[string]int)
		}
		commands[strings.Join(command.Command, " ")] += 1
	}
	commandsStat = addDailyCommands(commandsStat, commands, &date)

	lastCommandId := commandsHistory[len(commandsHistory)-1].ID
	return commandsStat, lastCommandId, nil
}

func buildInstancesStats(infra *graph.Graph) (instancesStats []*instancesStat, err error) {
	instancesStats, err = addStatsForInstanceStringProperty(infra, "Type", "InstanceType", instancesStats)
	if err != nil {
		return instancesStats, err
	}
	instancesStats, err = addStatsForInstanceStringProperty(infra, "ImageId", "ImageId", instancesStats)
	if err != nil {
		return instancesStats, err
	}

	return instancesStats, err
}

func addStatsForInstanceStringProperty(infra *graph.Graph, propertyName string, instanceStatType string, instancesStats []*instancesStat) ([]*instancesStat, error) {
	nodes, err := infra.NodesForType(graph.Instance.ToRDFString())
	if err != nil {
		return nil, err
	}
	propertyValuesCountMap := make(map[string]int)
	for _, i := range nodes {
		inst := graph.InitFromRdfNode(i)
		e := inst.UnmarshalFromGraph(infra)
		if e != nil {
			return nil, e
		}
		if inst.Properties()[propertyName] != nil {
			propertyValue, ok := inst.Properties()[propertyName].(string)
			if !ok {
				return nil, fmt.Errorf("Property value of '%s' is not a string: %T", inst.Properties()[propertyName], inst.Properties()[propertyName])
			}
			propertyValuesCountMap[propertyValue]++
		}
	}

	for k, v := range propertyValuesCountMap {
		instancesStats = append(instancesStats, &instancesStat{Type: instanceStatType, Date: time.Now(), Hits: v, Name: k})
	}

	return instancesStats, err
}

func addDailyCommands(commandsStat []*dailyCommands, commands map[string]int, date *time.Time) []*dailyCommands {
	for command, hits := range commands {
		dc := dailyCommands{Command: command, Hits: hits, Date: *date}
		commandsStat = append(commandsStat, &dc)
	}
	return commandsStat
}

type infraMetrics struct {
	Date                  time.Time
	Region                string
	NbVpcs                int
	NbSubnets             int
	MinSubnetsPerVpc      int
	MaxSubnetsPerVpc      int
	NbInstances           int
	MinInstancesPerSubnet int
	MaxInstancesPerSubnet int
}

func buildInfraMetrics(region string, infra *graph.Graph) (*infraMetrics, error) {
	metrics := &infraMetrics{
		Date:   time.Now(),
		Region: region,
	}

	c, min, max, err := computeCountMinMaxChildForType(infra, graph.Vpc)
	if err != nil {
		return metrics, err
	}
	metrics.NbVpcs, metrics.MinSubnetsPerVpc, metrics.MaxSubnetsPerVpc = c, min, max

	c, min, max, err = computeCountMinMaxChildForType(infra, graph.Subnet)
	if err != nil {
		return metrics, err
	}
	metrics.NbSubnets, metrics.MinInstancesPerSubnet, metrics.MaxInstancesPerSubnet = c, min, max

	c, _, _, err = computeCountMinMaxChildForType(infra, graph.Instance)
	if err != nil {
		return metrics, err
	}
	metrics.NbInstances = c

	return metrics, nil
}

func buildAccessMetrics(region string, access *graph.Graph, time time.Time) (*accessMetrics, error) {
	metrics := &accessMetrics{
		Date:   time,
		Region: region,
	}
	c, min, max, err := computeCountMinMaxForTypeWithChildType(access, graph.Group, graph.User)
	if err != nil {
		return metrics, err
	}
	metrics.NbGroups, metrics.MinUsersByGroup, metrics.MaxUsersByGroup = c, min, max

	c, min, max, err = computeCountMinMaxForTypeWithChildType(access, graph.Policy, graph.User)
	if err != nil {
		return metrics, err
	}
	metrics.NbPolicies, metrics.MinUsersByLocalPolicies, metrics.MaxUsersByLocalPolicies = c, min, max

	_, min, max, err = computeCountMinMaxForTypeWithChildType(access, graph.Policy, graph.Role)
	if err != nil {
		return metrics, err
	}
	metrics.MinRolesByLocalPolicies, metrics.MaxRolesByLocalPolicies = min, max

	_, min, max, err = computeCountMinMaxForTypeWithChildType(access, graph.Policy, graph.Group)
	if err != nil {
		return metrics, err
	}
	metrics.MinGroupsByLocalPolicies, metrics.MaxGroupsByLocalPolicies = min, max

	c, _, _, err = computeCountMinMaxChildForType(access, graph.Role)
	if err != nil {
		return metrics, err
	}
	metrics.NbRoles = c

	c, _, _, err = computeCountMinMaxChildForType(access, graph.User)
	if err != nil {
		return metrics, err
	}
	metrics.NbUsers = c

	return metrics, nil
}

func computeCountMinMaxChildForType(graph *graph.Graph, t graph.ResourceType) (int, int, int, error) {
	nodes, err := graph.NodesForType(t.ToRDFString())
	if err != nil {
		return 0, 0, 0, err
	}
	if len(nodes) == 0 {
		return 0, 0, 0, nil
	}
	firstNode := nodes[0]
	count, err := graph.CountChildrenForNode(firstNode)
	if err != nil {
		return 0, 0, 0, err
	}

	min, max := count, count
	for _, node := range nodes[1:] {
		count, err = graph.CountChildrenForNode(node)
		if err != nil {
			return 0, 0, 0, err
		}
		if count < min {
			min = count
		}
		if count > max {
			max = count
		}
	}
	return len(nodes), min, max, nil
}

func computeCountMinMaxForTypeWithChildType(graph *graph.Graph, parentType, childType graph.ResourceType) (int, int, int, error) {
	nodes, err := graph.NodesForType(parentType.ToRDFString())
	if err != nil {
		return 0, 0, 0, err
	}
	if len(nodes) == 0 {
		return 0, 0, 0, nil
	}
	firstNode := nodes[0]
	count, err := graph.CountChildrenOfTypeForNode(firstNode, childType)
	if err != nil {
		return 0, 0, 0, err
	}

	min, max := count, count
	for _, node := range nodes[1:] {
		count, err = graph.CountChildrenOfTypeForNode(node, childType)
		if err != nil {
			return 0, 0, 0, err
		}
		if count < min {
			min = count
		}
		if count > max {
			max = count
		}
	}
	return len(nodes), min, max, nil
}

func sameDay(date1, date2 *time.Time) bool {
	return (date1.Day() == date2.Day()) && (date1.Month() == date2.Month()) && (date1.Year() == date2.Year())
}

type encryptedData struct {
	Key  []byte
	Data []byte
}

func aesEncrypt(data []byte) ([]byte, []byte, error) {
	sessionKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, sessionKey); err != nil {
		return nil, nil, err
	}

	aesCipher, err := aes.NewCipher(sessionKey)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	encrypted := gcm.Seal(nonce, nonce, data, nil)
	return sessionKey, encrypted, nil
}
