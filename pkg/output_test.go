package pkg

import (
	"testing"
)

func TestSetFormat(t *testing.T) {
	SetFormat("json")
	if format != "json" {
		t.Fatalf(`failed to set format, wanted %s, got %s`, "json", format)
	}
}
