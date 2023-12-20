package base

import (
	"bytes"
	"testing"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	buf := &bytes.Buffer{}
	c := collection.New(buf)

	Collect(c)

	err := c.Close()
	assert.NoError(t, err)

	assert.NotEmpty(t, buf.Len())
}
