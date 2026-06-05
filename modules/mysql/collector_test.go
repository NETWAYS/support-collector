package mysql

import (
	"bytes"
	"testing"

	"github.com/NETWAYS/support-collector/internal/collection"
)

func TestCollect(t *testing.T) {
	c := collection.New(&bytes.Buffer{})

	Collect(c)

	err := c.Close()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
