package foreman

import (
	"bytes"
	"testing"

	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"github.com/NETWAYS/support-collector/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	c := collection.New(&bytes.Buffer{})

	if !detect() {
		t.Skip("could not find foreman in the test environment")
		return
	}

	Collect(c)

	err := c.Close()
	assert.NoError(t, err)
}

func TestObfuscators(t *testing.T) {
	util.AssertObfuscationExample(t, obfuscaters, obfuscate.KindFile, "/etc/foreman/database.yml")
	util.AssertObfuscationExample(t, obfuscaters, obfuscate.KindFile, "/etc/foreman/encryption_key.rb")

	util.AssertAllObfuscatorsTested(t, obfuscaters)
}
