package main

import (
	"fmt"
	"os"

	"github.com/tobischo/gokeepasslib/v3"
)

func main() {
	var customer string = "maz005"
	var environment string = "p"
	var prefix string = customer + "-" + environment
	// envvars := [5]string{"TF_VAR_tenant_id", "TF_VAR_client_id", "TF_VAR_client_secret", "TF_VAR_subscription_id", "ARM_ACCESS_KEY"}

	// for _, env := range envvars {
	// 	if env == "" {
	// 		fmt.Printf(os.Getenv(env))
	// 	} else {
	// 		fmt.Printf("%s not found.\n", env)
	// 	}
	// }

	file, _ := os.Open("")

	db := gokeepasslib.NewDatabase()
	db.Credentials = gokeepasslib.NewPasswordCredentials("")
	_ = gokeepasslib.NewDecoder(file).Decode(db)

	db.UnlockProtectedEntries()

	groups := db.Content.Root.Groups[0].Groups[0].Entries

	for _, group := range groups {
		// fmt.Println(group.GetTitle())
		if group.GetContent("Title") == prefix {
			fmt.Println(group.GetContent("Title"))
			fmt.Println(group.GetContent("UserName"))
			fmt.Println(group.GetContent("Password"))
		}
	}

}
