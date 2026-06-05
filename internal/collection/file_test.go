package collection

import (
	"bytes"
	"testing"
)

func TestLoadFiles(t *testing.T) {
	files, err := LoadFiles("test", "testdata/example.txt")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(files) != 1 {
		t.Error("Expected len(files) to be 1")
	}

	files, err = LoadFiles("test", "testdata")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(files) != 2 {
		t.Error("Expected len(files) to be 2")
	}

	files, err = LoadFiles("test", "testdata/*.txt")

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(files) != 1 {
		t.Error("Expected len(files) to be 1")
	}
}

func TestFile_Write(t *testing.T) {
	f := NewFile("test.txt")

	_, err := f.Write([]byte("content"))

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	_, err = f.Write([]byte("content"))

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !bytes.Equal(f.Data, []byte("contentcontent")) {
		t.Errorf("expected %q, got %q", []byte("contentcontent"), f.Data)
	}
}
