package ansible

import (
	"bytes"
	"testing"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/util"
)

func TestCollect(t *testing.T) {
	if !util.ModuleExists(relevantPaths) {
		t.Skip("could not find ansible in the test environment")
		return
	}

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
