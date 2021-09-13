package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/tobischo/gokeepasslib"
)

var (
	version    = "dev-build"
	goVersion  = runtime.Version()
	versionStr = fmt.Sprintf("CTX version %v, %v", version, goVersion)
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

func printForExport(lines []string) {
	for _, line := range lines {
		fmt.Printf(" export %v\n", line)
	}
}

func main() {
	environmentFlag := flag.String("e", "", "Specify the customer environment, which is the title of the Keepass secret (eg. maz000-p).")
	targetGroupNameFlag := flag.String("g", "Azure", "The Keepass group where the variables are stored.")
	versionFlag := flag.Bool("v", false, "Displays the version number of CTX and Go.")

	flag.Parse()

	if *versionFlag {
		fmt.Println(versionStr)
		os.Exit(0)
	}

	cfg := newConfig()

	file, _ := os.Open(cfg.dbLocation)
	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials(cfg.secret)
	_ = gokeepasslib.NewDecoder(file).Decode(db)
	db.UnlockProtectedEntries()

	groups := db.Content.Root.Groups[0].Groups
	group := findGroup(*targetGroupNameFlag, groups)

	notes := getNotes(group, *environmentFlag)
	printForExport(notes)

}
