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
	Id       string
	Version  string
	Commands []*DailyCommands
}

type DailyCommands struct {
	Command string
	Hits    int
	Date    time.Time
}

func (db *DB) BuildStats(fromCommandId int) (*Stats, int, error) {
	id, err := db.GetStringValue(AWLESS_ID_KEY)
	if err != nil {
		return nil, 0, err
	}

	stats := &Stats{Id: id, Version: config.Version, Commands: []*DailyCommands{}}
	commandsHistory, err := db.GetHistory(fromCommandId)
	if err != nil {
		return stats, 0, err
	}

	if len(commandsHistory) == 0 {
		return stats, 0, nil
	}

	date := commandsHistory[0].Time
	commands := make(map[string]int)

	for _, command := range commandsHistory {
		if !SameDay(&date, &command.Time) {
			addDailyCommands(stats, commands, &date)
			date = command.Time
			commands = make(map[string]int)
		}
		commands[strings.Join(command.Command, " ")] += 1
	}
	addDailyCommands(stats, commands, &date)

	lastCommandId := commandsHistory[len(commandsHistory)-1].Id
	return stats, lastCommandId, nil
}

func addDailyCommands(stats *Stats, commands map[string]int, date *time.Time) {
	for command, hits := range commands {
		dc := DailyCommands{Command: command, Hits: hits, Date: *date}
		stats.Commands = append(stats.Commands, &dc)
	}
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
	stats, lastCommandId, err := db.BuildStats(lastCommandId)
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
