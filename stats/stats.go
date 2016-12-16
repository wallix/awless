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

	"github.com/spf13/viper"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
)

const (
	AWLESS_ID_KEY  = "awless_id"
	AWLESS_AID_KEY = "awless_aid"
	SENT_ID_KEY    = "sent_id"
	SENT_TIME_KEY  = "sent_time"
)

type Stats struct {
	Id             string
	AId            string
	Version        string
	Commands       []*DailyCommands
	InfraMetrics   *InfraMetrics
	InstancesStats []*InstancesStat
	AccessMetrics  *AccessMetrics
	Logs           []*Log
}

type DailyCommands struct {
	Command string
	Hits    int
	Date    time.Time
}

type InstancesStat struct {
	Type string
	Date time.Time
	Hits int
	Name string
}

type AccessMetrics struct {
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

func BuildStats(db *DB, infra *rdf.Graph, access *rdf.Graph, fromCommandId int) (*Stats, int, error) {
	commandsStat, lastCommandId, err := buildCommandsStat(db, fromCommandId)
	if err != nil {
		return nil, 0, err
	}

	infraMetrics := &InfraMetrics{}
	if infra != nil {
		infraMetrics, err = buildInfraMetrics(infra)
		if err != nil {
			return nil, 0, err
		}
	}

	accessMetrics := &AccessMetrics{}
	if access != nil {
		accessMetrics, err = buildAccessMetrics(access, time.Now())
		if err != nil {
			return nil, 0, err
		}
	}

	instancesStats, err := buildInstancesStats(infra)
	if err != nil {
		return nil, 0, err
	}

	id, err := db.GetStringValue(AWLESS_ID_KEY)
	if err != nil {
		return nil, 0, err
	}

	aId, err := db.GetStringValue(AWLESS_AID_KEY)
	if err != nil {
		return nil, 0, err
	}

	logs, err := db.GetLogs()
	if err != nil {
		return nil, 0, err
	}

	stats := &Stats{
		Id:             id,
		AId:            aId,
		Version:        config.Version,
		Commands:       commandsStat,
		InfraMetrics:   infraMetrics,
		InstancesStats: instancesStats,
		AccessMetrics:  accessMetrics,
		Logs:           logs,
	}

	return stats, lastCommandId, nil
}

func buildCommandsStat(db *DB, fromCommandId int) ([]*DailyCommands, int, error) {
	var commandsStat []*DailyCommands

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
		if !SameDay(&date, &command.Time) {
			commandsStat = addDailyCommands(commandsStat, commands, &date)
			date = command.Time
			commands = make(map[string]int)
		}
		commands[strings.Join(command.Command, " ")] += 1
	}
	commandsStat = addDailyCommands(commandsStat, commands, &date)

	lastCommandId := commandsHistory[len(commandsHistory)-1].Id
	return commandsStat, lastCommandId, nil
}

func buildInstancesStats(infra *rdf.Graph) (instancesStats []*InstancesStat, err error) {
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

func addStatsForInstanceStringProperty(infra *rdf.Graph, propertyName string, instanceStatType string, instancesStats []*InstancesStat) ([]*InstancesStat, error) {
	nodes, err := infra.NodesForType(rdf.INSTANCE)
	if err != nil {
		return nil, err
	}
	propertyValuesCountMap := make(map[string]int)
	for _, inst := range nodes {
		var instProperties cloud.Properties
		instProperties, err = cloud.LoadPropertiesTriples(infra, inst)
		if err != nil {
			return nil, err
		}
		if instProperties[propertyName] != nil {
			propertyValue, ok := instProperties[propertyName].(string)
			if !ok {
				return nil, fmt.Errorf("Property value of '%s' is not a string: %T", instProperties[propertyName], instProperties[propertyName])
			}
			propertyValuesCountMap[propertyValue]++
		}
	}

	for k, v := range propertyValuesCountMap {
		instancesStats = append(instancesStats, &InstancesStat{Type: instanceStatType, Date: time.Now(), Hits: v, Name: k})
	}

	return instancesStats, err
}

func addDailyCommands(commandsStat []*DailyCommands, commands map[string]int, date *time.Time) []*DailyCommands {
	for command, hits := range commands {
		dc := DailyCommands{Command: command, Hits: hits, Date: *date}
		commandsStat = append(commandsStat, &dc)
	}
	return commandsStat
}

type InfraMetrics struct {
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

func buildInfraMetrics(infra *rdf.Graph) (*InfraMetrics, error) {
	metrics := &InfraMetrics{
		Date:   time.Now(),
		Region: viper.GetString("region"),
	}

	c, min, max, err := computeCountMinMaxChildForType(infra, rdf.VPC)
	if err != nil {
		return metrics, err
	}
	metrics.NbVpcs, metrics.MinSubnetsPerVpc, metrics.MaxSubnetsPerVpc = c, min, max

	c, min, max, err = computeCountMinMaxChildForType(infra, rdf.SUBNET)
	if err != nil {
		return metrics, err
	}
	metrics.NbSubnets, metrics.MinInstancesPerSubnet, metrics.MaxInstancesPerSubnet = c, min, max

	c, _, _, err = computeCountMinMaxChildForType(infra, rdf.INSTANCE)
	if err != nil {
		return metrics, err
	}
	metrics.NbInstances = c

	return metrics, nil
}

func buildAccessMetrics(access *rdf.Graph, time time.Time) (*AccessMetrics, error) {
	metrics := &AccessMetrics{
		Date:   time,
		Region: viper.GetString("region"),
	}
	c, min, max, err := computeCountMinMaxForTypeWithChildType(access, rdf.GROUP, rdf.USER)
	if err != nil {
		return metrics, err
	}
	metrics.NbGroups, metrics.MinUsersByGroup, metrics.MaxUsersByGroup = c, min, max

	c, min, max, err = computeCountMinMaxForTypeWithChildType(access, rdf.POLICY, rdf.USER)
	if err != nil {
		return metrics, err
	}
	metrics.NbPolicies, metrics.MinUsersByLocalPolicies, metrics.MaxUsersByLocalPolicies = c, min, max

	_, min, max, err = computeCountMinMaxForTypeWithChildType(access, rdf.POLICY, rdf.ROLE)
	if err != nil {
		return metrics, err
	}
	metrics.MinRolesByLocalPolicies, metrics.MaxRolesByLocalPolicies = min, max

	_, min, max, err = computeCountMinMaxForTypeWithChildType(access, rdf.POLICY, rdf.GROUP)
	if err != nil {
		return metrics, err
	}
	metrics.MinGroupsByLocalPolicies, metrics.MaxGroupsByLocalPolicies = min, max

	c, _, _, err = computeCountMinMaxChildForType(access, rdf.ROLE)
	if err != nil {
		return metrics, err
	}
	metrics.NbRoles = c

	c, _, _, err = computeCountMinMaxChildForType(access, rdf.USER)
	if err != nil {
		return metrics, err
	}
	metrics.NbUsers = c

	return metrics, nil
}

func computeCountMinMaxChildForType(graph *rdf.Graph, t string) (int, int, int, error) {
	nodes, err := graph.NodesForType(t)
	if err != nil {
		return 0, 0, 0, err
	}
	if len(nodes) == 0 {
		return 0, 0, 0, nil
	}
	firstNode := nodes[0]
	count, err := graph.CountTriplesForSubjectAndPredicate(firstNode, rdf.ParentOfPredicate)
	if err != nil {
		return 0, 0, 0, err
	}

	min, max := count, count
	for _, node := range nodes[1:] {
		count, err = graph.CountTriplesForSubjectAndPredicate(node, rdf.ParentOfPredicate)
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

func computeCountMinMaxForTypeWithChildType(graph *rdf.Graph, parentType, childType string) (int, int, int, error) {
	nodes, err := graph.NodesForType(parentType)
	if err != nil {
		return 0, 0, 0, err
	}
	if len(nodes) == 0 {
		return 0, 0, 0, nil
	}
	firstNode := nodes[0]
	count, err := graph.CountTriplesForSubjectAndPredicateObjectOfType(firstNode, rdf.ParentOfPredicate, childType)
	if err != nil {
		return 0, 0, 0, err
	}

	min, max := count, count
	for _, node := range nodes[1:] {
		count, err = graph.CountTriplesForSubjectAndPredicateObjectOfType(node, rdf.ParentOfPredicate, childType)
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

func SameDay(date1, date2 *time.Time) bool {
	return (date1.Day() == date2.Day()) && (date1.Month() == date2.Month()) && (date1.Year() == date2.Year())
}

type EncryptedData struct {
	Key  []byte
	Data []byte
}

func (db *DB) SendStats(url string, publicKey rsa.PublicKey, localInfra, localAccess *rdf.Graph) error {
	lastCommandId, err := db.GetIntValue(SENT_ID_KEY)
	if err != nil {
		return err
	}

	stats, lastCommandId, err := BuildStats(db, localInfra, localAccess, lastCommandId)
	if err != nil {
		return err
	}

	var zipped bytes.Buffer
	zippedW := gzip.NewWriter(&zipped)
	if err = json.NewEncoder(zippedW).Encode(stats); err != nil {
		return err
	}
	zippedW.Close()

	sessionKey, encrypted, err := aesEncrypt(zipped.Bytes())
	if err != nil {
		return err
	}
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, sessionKey, nil)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(EncryptedData{encryptedKey, encrypted})
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 2 * time.Second}
	if _, err := client.Post(url, "application/json", bytes.NewReader(payload)); err != nil {
		return err
	}

	if err := db.SetIntValue(SENT_ID_KEY, lastCommandId); err != nil {
		return err
	}
	if err := db.SetTimeValue(SENT_TIME_KEY, time.Now()); err != nil {
		return err
	}
	if err := db.FlushLogs(); err != nil {
		return err
	}
	return nil
}

func (db *DB) CheckStatsToSend(expirationDuration time.Duration) bool {
	sent, err := db.GetTimeValue(SENT_TIME_KEY)
	if err != nil {
		sent = time.Time{}
	}
	return (time.Since(sent) > expirationDuration)
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
