package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestSecretRetrieval(t *testing.T) {
	l := log.New(os.Stdout, "test", log.LstdFlags)
	config := ctxConfig{
		environment: "maz998-t",
		groupName:   "TestGroup",
	}
	t.Setenv("CTX_VAR_secret", "Testpassword123!")
	t.Setenv("CTX_VAR_db_location", "/Users/jvanhienen/iCloud/Projects/ctx/mock/mock.kdbx")
	_, err := handleKeepass(l, config)

	if err != nil {
		t.Fail()
	}
}

// func TestSecretFormatting(t *testing.T) {
// 	l := log.New(os.Stdout, "test", log.LstdFlags)
// 	config := ctxConfig{
// 		environment: "maz998-t",
// 		groupName:   "TestGroup",
// 	}
// 	t.Setenv("CTX_VAR_secret", "Testpassword123!")
// 	t.Setenv("CTX_VAR_db_location", "/Users/jvanhienen/iCloud/Projects/ctx/mock/mock.kdbx")
//
// 	secrets, _ := handleKeepass(l, config)
// 	formattedSecrets := printForExport(secrets)
// 	if strings.HasPrefix(formattedSecrets[0], " export=") == false {
// 		t.Fail()
// 	}
// }

func TestInvalidDbLocation(t *testing.T) {
	l := log.New(os.Stdout, "test", log.LstdFlags)
	config := ctxConfig{
		environment: "maz998-t",
		groupName:   "TestGroup",
	}
	t.Setenv("CTX_VAR_db_location", "/tmp/wrong.kdbx")
	t.Setenv("CTX_VAR_secret", "Testpassword123!")

	_, err := handleKeepass(l, config)
	expectedErrorMsg := "open /tmp/wrong.kdbx: no such file or directory"
	if err != nil {
		assert.EqualError(t, err, expectedErrorMsg)
	}
	if err == nil {
		t.Errorf("Did not throw an error. got: nil, want: %s", expectedErrorMsg)
	}
}
