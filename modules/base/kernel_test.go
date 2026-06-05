package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKernelInfo(t *testing.T) {
	k, err := GetKernelInfo()

	assert.NoError(t, err)
	assert.NotEmpty(t, k.Kernel)
	assert.NotEmpty(t, k.Version)
	assert.NotEmpty(t, k.Hostname)
	assert.NotEmpty(t, k.FQDN)
}
