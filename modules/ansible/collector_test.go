package ansible

import (
	"bytes"
	"github.com/NETWAYS/support-collector/internal/util"
	"testing"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	assert.NotEmpty(t, buf.Len())
}
