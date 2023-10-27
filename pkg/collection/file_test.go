package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFiles(t *testing.T) {
	files, err := LoadFiles("test", "testdata/example.txt")
	assert.NoError(t, err)
	assert.Len(t, files, 1)

	files, err = LoadFiles("test", "testdata")
	assert.NoError(t, err)
	assert.Len(t, files, 2)

	files, err = LoadFiles("test", "testdata/*.txt")
	assert.NoError(t, err)
	assert.Len(t, files, 1)
}

func TestFile_Write(t *testing.T) {
	f := NewFile("test.txt")

	_, err := f.Write([]byte("content"))
	assert.NoError(t, err)

	_, err = f.Write([]byte("content"))
	assert.NoError(t, err)

	assert.Equal(t, f.Data, []byte("contentcontent"))
}
