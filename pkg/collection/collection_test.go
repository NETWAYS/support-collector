package collection

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollection_AddFileFromReader(t *testing.T) {
	c := &Collection{}

	err := c.AddFileFromReader("test.txt", bytes.NewBufferString("content"))
	assert.NoError(t, err)
	assert.Len(t, c.Files, 1)
}

func TestCollection_AddFiles(t *testing.T) {
	c := &Collection{}

	err := c.AddFiles("test", "testdata/")
	assert.NoError(t, err)
	assert.Len(t, c.Files, 2)
}
