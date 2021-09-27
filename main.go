package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/johnvanhienen/ctx/data"
)

var (
	version    = "dev-build"
	goVersion  = runtime.Version()
	versionStr = fmt.Sprintf("CTX version %v, %v", version, goVersion)
)

func printForExport(lines []string) {
	for _, line := range lines {
		fmt.Printf(" export %v\n", line)
	}
}

func main() {
	l := log.New(os.Stdout, "ctx", log.LstdFlags)
	environmentFlag := flag.String("e", "", "Specify the customer environment, which is the title of the Keepass secret (eg. maz000-p).")
	targetGroupNameFlag := flag.String("g", "Azure", "The Keepass group where the variables are stored.")
	dataSourceFlag := flag.String("d", "keepass", "Specify the data source for the secrets (eg. keyvault or keepass")
	versionFlag := flag.Bool("v", false, "Displays the version number of CTX and Go.")

	flag.Parse()

	if *versionFlag {
		fmt.Println(versionStr)
		os.Exit(0)
	}
	if *dataSourceFlag == "keepass" {
		kpcfg := data.NewKeepassConfig()
		kp := data.NewKeepass(l, kpcfg)
		kp.GroupName = *targetGroupNameFlag
		kp.Environment = *environmentFlag
		// fmt.Println(kp.GroupName)
		secrets := kp.GetSecrets()
		for _, secret := range secrets {
			fmt.Println(secret)
		}
	}
	// cfg := newConfig()

	// db, err := InitiliazeDatabase(cfg.dbLocation, cfg.secret)
	// if err != nil {
	// 	fmt.Println("Failed initializing database: ", err)
	// }

	// groups := db.Content.Root.Groups[0].Groups
	// group := findGroup(*targetGroupNameFlag, groups)

	// notes := getNotes(group, *environmentFlag)
	// printForExport(notes)

}
