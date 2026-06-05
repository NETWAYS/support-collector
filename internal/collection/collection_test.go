package collection

import (
	"bytes"
	"testing"
)

func TestCollection_AddFileFromReader(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(buf)

	err := c.AddFileFromReaderRaw("test.txt", bytes.NewBufferString("content"))

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	err = c.Close()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Expected buffer to be not empty")
	}
}

func TestCollection_AddFiles(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(buf)

	c.AddFiles("test", "testdata/")

	err := c.Close()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Expected buffer to be not empty")
	}
}
