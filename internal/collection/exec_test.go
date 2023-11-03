package collection

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testTimeout = 100 * time.Millisecond

func TestLoadCommandOutputWithTimeout(t *testing.T) {
	output, err := LoadCommandOutputWithTimeout(testTimeout, "sh", "-c", "echo good; echo stderr >&2")
	assert.NoError(t, err)
	assert.Equal(t, []byte("good\nstderr\n"), output)

	output, err = LoadCommandOutputWithTimeout(testTimeout, "sh", "-c", "exit 1")
	assert.Error(t, err)
	assert.NotEmpty(t, output)

	output, err = LoadCommandOutputWithTimeout(testTimeout, "sh", "-c", "sleep 1")
	assert.Error(t, err)
	assert.NotEmpty(t, output)
}
