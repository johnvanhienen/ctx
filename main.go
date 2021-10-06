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

type ctxConfig struct {
	environment string
	groupName   string
	dataSource  string
}

func main() {
	l := log.New(os.Stdout, "ctx", log.LstdFlags)
	environmentFlag := flag.String("e", "", "Specify the customer environment, which is the title of the Keepass secret (eg. maz000-p).")
	groupNameFlag := flag.String("g", "Azure", "The Keepass group where the variables are stored.")
	dataSourceFlag := flag.String("d", "keepass", "Specify the data source for the secrets (eg. keyvault or keepass")
	versionFlag := flag.Bool("v", false, "Displays the version number of CTX and Go.")
	flag.Parse()
	cfg := ctxConfig{
		environment: *environmentFlag,
		groupName:   *groupNameFlag,
		dataSource:  *dataSourceFlag,
	}

	if *versionFlag {
		fmt.Println(versionStr)
		os.Exit(0)
	}

	if cfg.dataSource == "keepass" {
		kpSecrets, err := handleKeepass(l, cfg)
		if err != nil {
			l.Fatal(err)
		}
		printForExport(kpSecrets)
	}
}

func handleKeepass(l *log.Logger, ctxCfg ctxConfig) ([]string, error) {
	kpCfg := data.NewKeepassConfig()
	kp := data.NewKeepass(l, kpCfg)
	kp.GroupName = ctxCfg.groupName
	kp.Environment = ctxCfg.environment
	result, err := kp.GetSecrets()
	return result, err

}

func printForExport(lines []string) {
	for _, line := range lines {
		fmt.Printf(" export %v\n", line)
	}
}
