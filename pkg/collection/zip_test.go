package collection

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestCollection_WriteZIP(t *testing.T) {
	tmp, err := ioutil.TempFile(os.TempDir(), "support-collector*.zip")
	assert.NoError(t, err)

	defer os.Remove(tmp.Name())

	c := &Collection{}

	assert.NoError(t, c.AddFile("test.txt", bytes.NewBufferString("content")))
	assert.NoError(t, c.AddFile("path/test2.txt", bytes.NewBufferString("content2")))

	assert.NoError(t, c.WriteZIP(tmp))

	fi, err := tmp.Stat()
	assert.NoError(t, err)
	assert.Greater(t, fi.Size(), int64(0))
}
