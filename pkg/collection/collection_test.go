package collection

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollection_AddFile(t *testing.T) {
	c := &Collection{}

	err := c.AddFile("test.txt", bytes.NewBufferString("content"))
	assert.NoError(t, err)
	assert.Len(t, c.Files, 1)
}
