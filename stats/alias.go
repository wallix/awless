package stats

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

type Aliases map[string]string

const (
	ALIASES_KEY = "aliases"
)

func (db *DB) GetAliases() (Aliases, error) {
	aliases := make(Aliases)
	b, err := db.GetValue(ALIASES_KEY)
	if err != nil {
		return aliases, err
	}
	if len(b) == 0 {
		return aliases, nil
	}
	err = json.Unmarshal(b, &aliases)
	return aliases, err
}

func (db *DB) AddAlias(name, target string) error {
	aliases, err := db.GetAliases()
	if err != nil {
		return err
	}
	name = strings.TrimSpace(name)
	target = strings.TrimSpace(target)
	if containsSpace(name) || containsSpace(target) {
		return fmt.Errorf("An alias must not contain any space")
	}
	if aliases[target] != "" {
		target = aliases[target]
	}
	if name == target {
		return fmt.Errorf("Useless alias (name and target are identical)")
	}
	if name == "" || target == "" {
		return fmt.Errorf("Alias name and target must not be empty")
	}

	aliases[name] = target
	return saveAliases(db, aliases)
}

func (db *DB) DeleteAlias(names ...string) error {
	aliases, err := db.GetAliases()
	if err != nil {
		return err
	}
	for _, name := range names {
		name = strings.TrimSpace(name)
		delete(aliases, name)
	}
	return saveAliases(db, aliases)
}

func saveAliases(db *DB, aliases Aliases) error {
	b, err := json.Marshal(aliases)
	if err != nil {
		return err
	}
	return db.SetValue(ALIASES_KEY, b)
}

func containsSpace(str string) bool {
	for _, ch := range str {
		if unicode.IsSpace(ch) {
			return true
		}
	}
	return false
}
