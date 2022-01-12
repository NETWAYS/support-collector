package ansible

import (
	"bytes"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollect(t *testing.T) {
	buf := &bytes.Buffer{}
	c := collection.New(buf)

	Collect(c)

	err := c.Close()
	assert.NoError(t, err)

	assert.NotEmpty(t, buf.Len())
}
