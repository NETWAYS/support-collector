package foreman

import (
	"bytes"
	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/collection/obfuscate"
	"github.com/NETWAYS/support-collector/internal/util"
	"testing"
)

func TestCollect(t *testing.T) {
	c := collection.New(&bytes.Buffer{})

	if !detect() {
		t.Skip("could not find foreman in the test environment")
		return
	}

	Collect(c)
}

func TestObfuscators(t *testing.T) {
	util.AssertObfuscationExample(t, obfuscaters, obfuscate.KindFile, "/etc/foreman/database.yml")
	util.AssertObfuscationExample(t, obfuscaters, obfuscate.KindFile, "/etc/foreman/encryption_key.rb")

	util.AssertAllObfuscatorsTested(t, obfuscaters)
}
