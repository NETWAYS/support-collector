package graphite

import (
	"bytes"
	"github.com/NETWAYS/support-collector/internal/util"
	"testing"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	if !util.ModuleExists(relevantPaths) {
		t.Skip("could not find graphite in the test environment")
		return
	}

	c := collection.New(&bytes.Buffer{})

	Collect(c)

	err := c.Close()
	assert.NoError(t, err)
}
