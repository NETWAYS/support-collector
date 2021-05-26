package collection

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestCollection_WriteZIP(t *testing.T) {
	tmp, err := ioutil.TempFile(os.TempDir(), "support-collector*.zip")
	assert.NoError(t, err)

	defer os.Remove(tmp.Name())

	c := New()

	assert.NoError(t, c.AddFiles("test", "testdata/"))

	assert.NoError(t, c.WriteZIP(tmp))

	fi, err := tmp.Stat()
	assert.NoError(t, err)
	assert.Greater(t, fi.Size(), int64(0))
}
