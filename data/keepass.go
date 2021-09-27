package data

import (
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
	c.location = os.Getenv("HOME") + "/database.kdbx"

	if os.Getenv("CTX_VAR_db_location") != "" {
		c.location = os.Getenv("CTX_VAR_db_location")
	}

	if os.Getenv("CTX_VAR_secret") == "" {
		fmt.Println("Set Keepass secret with CTX_VAR_secret")
		os.Exit(1)
	} else {
		c.secret = os.Getenv("CTX_VAR_secret")
	}

	return &c
}

func findGroup(targetGroupName string, groupPoolPtr []gokeepasslib.Group) (result gokeepasslib.Group) {
	for _, group := range groupPoolPtr {
		if group.Name == targetGroupName {
			result := group
			return result
		}
	}
	fmt.Printf("Root group %s was not found.", targetGroupName)
	return result
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

func initDatabase(c *keepassConfig) (*gokeepasslib.Database, error) {
	databaseFile, err := os.Open(c.location)
	if err != nil {
		fmt.Println("Can't find the database: ", err)
	}

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(c.secret)
	err = gokeepasslib.NewDecoder(databaseFile).Decode(db)

	if err != nil {
		fmt.Println("Can't open the database: ", err)
		return nil, err
	}

	db.UnlockProtectedEntries()

	return db, nil
}

func (s *Secrets) GetSecrets() []string {
	db, err := initDatabase(s.dbConfig)
	if err != nil {
		s.l.Fatal(err)
	}
	allGroups := db.Content.Root.Groups[0].Groups
	group := findGroup(s.GroupName, allGroups)
	secrets := getNotes(group, s.Environment)
	// if secrets != nil {
	// 	fmt.Println()
	// }
	// fmt.Println(s.GroupName)
	return secrets

}
