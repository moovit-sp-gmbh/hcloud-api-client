package config

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	identifier := "test" + randomString(8)
	LoadConfig()
	AddContext(identifier, "https://app.helmut.cloud", "test@mail.de", "testToken")

	if Config.GetActiveContext() != nil && Config.GetActiveContext().Identifier == identifier {
		t.Fatalf("failed to delete context, wanted not %s, got %s", "test", Config.GetActiveContext().Identifier)
	}

	SetContext(identifier)

	if Config.GetActiveContext().Identifier != identifier {
		t.Fatalf("failed to set context, wanted %s, got %s", identifier, Config.GetActiveContext().Identifier)
	}

	DelContext(identifier)

	if Config.GetActiveContext() != nil && Config.GetActiveContext().Identifier == identifier {
		t.Fatalf("failed to delete context, wanted not %s, got %s", identifier, Config.GetActiveContext().Identifier)
	}
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	var alphabet string = "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder

	l := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(l)]
		sb.WriteByte(c)
	}

	return sb.String()
}
