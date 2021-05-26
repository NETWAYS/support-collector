package collection

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile_Write(t *testing.T) {
	f := NewFile("test.txt")

	_, err := f.Write([]byte("content"))
	assert.NoError(t, err)

	_, err = f.Write([]byte("content"))
	assert.NoError(t, err)

	assert.Equal(t, f.Data, []byte("contentcontent"))
}
