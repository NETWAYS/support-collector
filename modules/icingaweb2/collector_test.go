package icingaweb2

import (
	"bytes"
	"testing"

	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"github.com/NETWAYS/support-collector/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	if !Detect() {
		t.Skip("could not find icingaweb2 in the test environment")
		return
	}

	c := collection.New(&bytes.Buffer{})

	Collect(c)

	err := c.Close()
	assert.NoError(t, err)
}

func TestObfuscators(t *testing.T) {
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/etc/icingaweb2/resources.ini")

	util.AssertAllObfuscatorsTested(t, obfuscators)
}
