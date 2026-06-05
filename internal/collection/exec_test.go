package collection

import (
	"bytes"
	"testing"
	"time"
)

var testTimeout = 100 * time.Millisecond

func TestLoadCommandOutputWithTimeout(t *testing.T) {
	output, err := LoadCommandOutputWithTimeout(testTimeout, "sh", "-c", "echo good; echo stderr >&2")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !bytes.Equal(output, []byte("good\nstderr\n")) {
		t.Errorf("expected %q, got %q", []byte("good\nstderr\n"), output)
	}

	output, err = LoadCommandOutputWithTimeout(testTimeout, "sh", "-c", "exit 1")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if bytes.Equal(output, []byte{}) {
		t.Errorf("expected output not to be empty")
	}

	output, err = LoadCommandOutputWithTimeout(testTimeout, "sh", "-c", "sleep 1")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if bytes.Equal(output, []byte{}) {
		t.Errorf("expected output not to be empty")
	}
}
