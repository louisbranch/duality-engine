package scenario

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestAssertionsLogOnly(t *testing.T) {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "", 0)

	assertions := Assertions{
		Mode:   AssertionLogOnly,
		Logger: logger,
	}

	if err := assertions.Assertf("expected %s", "thing"); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if !strings.Contains(buffer.String(), "expected thing") {
		t.Fatalf("expected log output to include assertion message")
	}
}
