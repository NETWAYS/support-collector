package base

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetKernelInfo(t *testing.T) {
	k, err := GetKernelInfo()

	assert.NoError(t, err)
	assert.NotEmpty(t, k.Kernel)
	assert.NotEmpty(t, k.Version)
	assert.NotEmpty(t, k.Hostname)
	assert.NotEmpty(t, k.FQDN)
}
