package data

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tobischo/gokeepasslib"
)

type Secrets struct {
	l           *log.Logger
	GroupName   string
	Environment string
	dbConfig    *keepassConfig
}

type keepassConfig struct {
	location string
	secret   string
}

func NewKeepass(l *log.Logger, c *keepassConfig) *Secrets {
	return &Secrets{l, "", "", c}
}

func NewKeepassConfig() *keepassConfig {
	c := keepassConfig{}

	if os.Getenv("CTX_VAR_db_location") != "" {
		c.location = os.Getenv("CTX_VAR_db_location")
	} else {
		c.location = os.Getenv("HOME") + "/database.kdbx"
	}

	if os.Getenv("CTX_VAR_secret") == "" {
		fmt.Println("Set Keepass secret with CTX_VAR_secret")
		os.Exit(1)
	} else {
		c.secret = os.Getenv("CTX_VAR_secret")
	}

	return &c
}

func findGroup(targetGroupName string, groupPoolPtr []gokeepasslib.Group) (result gokeepasslib.Group, err error) {
	if targetGroupName == "" {
		return result, errors.New("root group name missing")
	}

	for _, group := range groupPoolPtr {
		if group.Name == targetGroupName {
			result := group
			return result, nil
		}
	}
	return result, fmt.Errorf("root group %s was not found", targetGroupName)
}

func getNotes(group gokeepasslib.Group, environment string) (result []string) {
	for _, entry := range group.Entries {
		if entry.GetContent("Title") == environment {
			notes := entry.GetContent("Notes")
			result = strings.Split(notes, "\n")
		}
	}

	return result
}

func initDatabase(c *keepassConfig) (db *gokeepasslib.Database, err error) {
	databaseFile, err := os.Open(c.location)
	if err != nil {
		return nil, err
	}

	db = gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(c.secret)
	err = gokeepasslib.NewDecoder(databaseFile).Decode(db)
	if err != nil {
		return nil, err
	}

	err = db.UnlockProtectedEntries()
	if err != nil {
		return nil, err
	}

	return db, err
}

func (s *Secrets) GetSecrets() ([]string, error) {
	db, err := initDatabase(s.dbConfig)
	if err != nil {
		return nil, err
	}
	allGroups := db.Content.Root.Groups[0].Groups
	group, err := findGroup(s.GroupName, allGroups)
	if err != nil {
		return nil, err
	}
	secrets := getNotes(group, s.Environment)

	return secrets, nil
}
