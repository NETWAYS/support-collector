package mysql

import (
	"bytes"
	"testing"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	c := collection.New(&bytes.Buffer{})

	Collect(c)

	err := c.Close()
	assert.NoError(t, err)
}
