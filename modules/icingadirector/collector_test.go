package icingadirector

import (
	"bytes"
	"testing"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/util"
)

func TestCollect(t *testing.T) {
	if !util.ModuleExists([]string{InstallationPath}) {
		t.Skip("could not find icingadirector in the test environment")
		return
	}

	c := collection.New(&bytes.Buffer{})

	Collect(c)

	err := c.Close()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
