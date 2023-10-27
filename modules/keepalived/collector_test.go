package keepalived

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

	Collect(c)

	err := c.Close()
	assert.NoError(t, err)
}

func TestObfuscators(t *testing.T) {
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/etc/keepalived/keepalived.conf")

	util.AssertAllObfuscatorsTested(t, obfuscators)
}
