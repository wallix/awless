package stats

import (
	"bytes"
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
	Id         string
	DailyStats []*DailyStat
}

type DailyStat struct {
	Commands map[string]int
	Date     time.Time
}

func (db *DB) BuildStats(fromCommandId int) (*Stats, int, error) {
	id, err := db.GetStringValue(AWLESS_ID_KEY)
	if err != nil {
		return nil, 0, err
	}

	stats := &Stats{Id: id, DailyStats: []*DailyStat{}}
	commands, err := db.GetHistory(fromCommandId)
	if err != nil {
		return stats, 0, err
	}

	if len(commands) == 0 {
		return stats, 0, nil
	}

	dailyStat := &DailyStat{make(map[string]int), commands[0].Time}
	lastCommandId := commands[0].Id

	for _, command := range commands {
		if !SameDay(&dailyStat.Date, &command.Time) {
			stats.DailyStats = append(stats.DailyStats, dailyStat)
			dailyStat = &DailyStat{make(map[string]int), command.Time}
		}
		dailyStat.Commands[strings.Join(command.Command, " ")] += 1
		lastCommandId = command.Id
	}

	stats.DailyStats = append(stats.DailyStats, dailyStat)

	return stats, lastCommandId, nil
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
	marshaled, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	sessionKey, encrypted, err := aesEncrypt(marshaled)
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
