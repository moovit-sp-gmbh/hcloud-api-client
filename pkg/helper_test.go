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
	testCode := "000.000.0000"
	testError := "test.error"
	testMessage := "test message"
	body := fmt.Sprintf(`{"code":"%s","error":"%s","message":"%s"}`, testCode, testError, testMessage)
	erro := ParseError(body)
	if erro.Code != "000.000.0000" {
		t.Fatalf(`result didnt match, wanted %s, got %s`, testCode, erro.Code)
	}
	if erro.Error != testError {
		t.Fatalf(`result didnt match, wanted %s, got %s`, testError, erro.Error)
	}
	if erro.Message != testMessage {
		t.Fatalf(`result didnt match, wanted %s, got %s`, testMessage, erro.Message)
	}
}
