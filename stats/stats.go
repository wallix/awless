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
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
)

const AWLESS_ID_KEY = "awless_id"
const SENT_ID_KEY = "sent_id"
const SENT_TIME_KEY = "sent_time"

func generateAwlessId() (string, error) {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(seed)), nil
}

type Stats struct {
	Id             string
	Version        string
	Commands       []*DailyCommands
	InfraMetrics   *InfraMetrics
	InstancesStats []*InstancesStat
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

func BuildStats(db *DB, infra *rdf.Graph, fromCommandId int) (*Stats, int, error) {

	commandsStat, lastCommandId, err := buildCommandsStat(db, fromCommandId)
	if err != nil {
		return nil, 0, err
	}

	infraMetrics := &InfraMetrics{}
	if infra != nil {
		infraMetrics, err = buildInfraMetrics(infra, time.Now())
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

	stats := &Stats{
		Id:             id,
		Version:        config.Version,
		Commands:       commandsStat,
		InfraMetrics:   infraMetrics,
		InstancesStats: instancesStats,
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
	triples, err := infra.TriplesForPredicateName("Type")
	if err != nil {
		return instancesStats, err
	}

	instanceTypes := make(map[string]int)
	for _, t := range triples {
		l, e := t.Object().Literal()
		if e != nil {
			return instancesStats, e
		}
		instanceType, e := l.Text()
		if e != nil {
			return instancesStats, e
		}
		instanceTypes[instanceType]++
	}

	for k, v := range instanceTypes {
		instancesStats = append(instancesStats, &InstancesStat{Type: "InstanceType", Date: time.Now(), Hits: v, Name: k})
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

func buildInfraMetrics(infra *rdf.Graph, time time.Time) (*InfraMetrics, error) {
	metrics := &InfraMetrics{
		Date:   time,
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

func computeCountMinMaxChildForType(graph *rdf.Graph, t string) (int, int, int, error) {
	nodes, err := graph.NodesForType(t)
	if err != nil {
		return 0, 0, 0, err
	}
	if len(nodes) == 0 {
		return 0, 0, 0, nil
	}
	firstNode := nodes[0]
	count, err := graph.CountTriplesForSubjectAndPredicate(firstNode, rdf.ParentOf)
	if err != nil {
		return 0, 0, 0, err
	}

	min, max := count, count
	for _, node := range nodes[1:] {
		count, err = graph.CountTriplesForSubjectAndPredicate(node, rdf.ParentOf)
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

func (db *DB) SendStats(url string, publicKey rsa.PublicKey) error {
	lastCommandId, err := db.GetIntValue(SENT_ID_KEY)
	if err != nil {
		return err
	}

	localInfra, err := rdf.NewGraphFromFile(filepath.Join(config.Dir, config.InfraFilename))
	if err != nil {
		return err
	}

	stats, lastCommandId, err := BuildStats(db, localInfra, lastCommandId)
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
