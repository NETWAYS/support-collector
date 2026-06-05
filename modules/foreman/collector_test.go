package foreman

import (
	"bytes"
	"testing"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
	"github.com/NETWAYS/support-collector/internal/util"
)

func TestCollect(t *testing.T) {
	if !util.ModuleExists(relevantPaths) {
		t.Skip("could not find foreman in the test environment")
		return
	}

	c := collection.New(&bytes.Buffer{})

	Collect(c)

	err := c.Close()

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestObfuscators(t *testing.T) {
	util.AssertObfuscationExample(t, obfuscaters, obfuscate.KindFile, "/etc/foreman/database.yml")
	util.AssertObfuscationExample(t, obfuscaters, obfuscate.KindFile, "/etc/foreman/encryption_key.rb")

	util.AssertAllObfuscatorsTested(t, obfuscaters)
}
