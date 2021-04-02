package main

import (
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
	var environment string = "maz005-p"
	// envvars := [5]string{"TF_VAR_tenant_id", "TF_VAR_client_id", "TF_VAR_client_secret", "TF_VAR_subscription_id", "ARM_ACCESS_KEY"}

	// for _, env := range envvars {
	// 	if env == "" {
	// 		fmt.Printf(os.Getenv(env))
	// 	} else {
	// 		fmt.Printf("%s not found.\n", env)
	// 	}
	// }

	cfg := newConfig()

	file, _ := os.Open(cfg.dbLocation)

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(cfg.secret)
	_ = gokeepasslib.NewDecoder(file).Decode(db)

	db.UnlockProtectedEntries()

	groups := db.Content.Root.Groups[0].Groups[0].Entries

	for _, group := range groups {
		// fmt.Println(group.GetTitle())
		if group.GetContent("Title") == environment {
			fmt.Println(group.GetContent("Title"))
			fmt.Println(group.GetContent("UserName"))
			fmt.Println(group.GetContent("Password"))
		}
	}

}
