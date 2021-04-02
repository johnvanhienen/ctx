package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tobischo/gokeepasslib"
)

type config struct {
	dbLocation string
	secret     string
}

func newConfig() *config {
	c := config{}
	c.dbLocation = os.Getenv("HOME") + "/database.kdbx"

	if os.Getenv("CTX_VAR_db_location") != "" {
		c.dbLocation = os.Getenv("CTX_VAR_db_location")
	}

	if os.Getenv("CTX_VAR_secret") == "" {
		fmt.Println("Set Keepass secret with CTX_VAR_secret")
		os.Exit(1)
	} else {
		c.secret = os.Getenv("CTX_VAR_secret")
	}

	return &c
}

func main() {
	environmentPtr := flag.String("e", "", "Specify the customer environment, which is the title of the Keepass secret (eg. maz000-p).")
	keepassGroupPtr := flag.String("g", "Azure", "The Keepass group where the variables are stored.")
	flag.Parse()

	cfg := newConfig()

	file, _ := os.Open(cfg.dbLocation)
	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(cfg.secret)
	_ = gokeepasslib.NewDecoder(file).Decode(db)
	db.UnlockProtectedEntries()

	var groupFound bool = false

	groups := db.Content.Root.Groups[0].Groups
	for _, group := range groups {
		if group.Name == *keepassGroupPtr {
			groupFound = true
			for _, entry := range group.Entries {
				if entry.GetContent("Title") == *environmentPtr {
					username := entry.GetContent("UserName")
					password := entry.GetContent("Password")
					fmt.Printf("export %s=%v\n", username, password)
				}
			}
		}
	}

	if !groupFound {
		fmt.Printf("Root group %s was not found.", *keepassGroupPtr)
	}

}
