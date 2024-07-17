package auth

import (
	"fmt"
	"os"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	dir, err := os.UserHomeDir()
	if err != nil {
		t.FailNow()
	}

	expectedString := fmt.Sprintf("%s/.texttunnel/config.json", dir)
	actualString, actualErr := getConfigPath()

	if expectedString != actualString {
		t.Fatalf("Expected %s, Received %s", expectedString, actualString)
	}

	if actualErr != nil {
		t.Fatalf("Expected nil error, received %s", actualErr)
	}
}

func TestSaveToken(t *testing.T) {
	err := SaveToken("TestToken")
	if err != nil {
		t.Fatalf("Expected nil, received err: %s", err)
	}
}
