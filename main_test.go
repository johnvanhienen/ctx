package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"
)

func TestSecretRetrieval(t *testing.T) {
	l := log.New(os.Stdout, "test", log.LstdFlags)
	config := ctxConfig{
		environment: "maz998-t",
		groupName:   "Azure",
	}
	_, err := handleKeepass(l, config)

	if err != nil {
		t.Fail()
	}
}

func TestSecretFormatting(t *testing.T) {
	l := log.New(os.Stdout, "test", log.LstdFlags)
	config := ctxConfig{
		environment: "maz998-t",
		groupName:   "Azure",
	}
	secrets, _ := handleKeepass(l, config)
	formattedSecrets := printForExport(secrets)
	if strings.HasPrefix(formattedSecrets[0], " export=") == false {
		t.Fail()
	}
}

func TestInvalidDbLocation(t *testing.T) {
	l := log.New(os.Stdout, "test", log.LstdFlags)
	config := ctxConfig{
		environment: "maz998-t",
		groupName:   "Azure",
	}
	t.Setenv("CTX_VAR_db_location", "/tmp/wrong.kdbx")

	_, err := handleKeepass(l, config)
	expectedErrorMsg := "open /tmp/wrong.kdbx: no such file or directory"
	if err != nil {
		assert.EqualError(t, err, expectedErrorMsg)
	}
	if err == nil {
		t.Errorf("Did not throw an error. got: nil, want: %s", expectedErrorMsg)
	}
}
