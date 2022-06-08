package pkg

import (
	"fmt"
	"os"
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	result := GetHomeDir()
	homeDir, _ := os.UserHomeDir()
	if result != homeDir {
		t.Fatalf(`result didnt match, wanted %s, got %s`, homeDir, result)
	}
}

func TestParseErr(t *testing.T) {
	testCode := -1
	testMessage := "test message"
	erro := ParseError(fmt.Sprintf(`{"status":%d,"message":"%s"}`, testCode, testMessage))
	if erro.Code != -1 {
		t.Fatalf(`result didnt match, wanted %d, got %d`, testCode, erro.Code)
	}
	if erro.Message != testMessage {
		t.Fatalf(`result didnt match, wanted %s, got %s`, testMessage, erro.Message)
	}
}
