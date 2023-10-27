package collection

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollection_AddFileFromReader(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(buf)

	err := c.AddFileFromReaderRaw("test.txt", bytes.NewBufferString("content"))
	assert.NoError(t, err)

	err = c.Close()
	assert.NoError(t, err)

	assert.Greater(t, buf.Len(), 0)
}

func TestCollection_AddFiles(t *testing.T) {
	buf := &bytes.Buffer{}
	c := New(buf)

	c.AddFiles("test", "testdata/")

	err := c.Close()
	assert.NoError(t, err)

	assert.Greater(t, buf.Len(), 0)
}
