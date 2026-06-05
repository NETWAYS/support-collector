package base

import (
	"bytes"
	"testing"

	"github.com/NETWAYS/support-collector/internal/collection"
)

func TestCollect(t *testing.T) {
	buf := &bytes.Buffer{}
	c := collection.New(buf)

	Collect(c)

	err := c.Close()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Expected buffer to be not empty")
	}
}
