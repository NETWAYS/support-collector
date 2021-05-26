package base

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetKernelInfo(t *testing.T) {
	k := GetKernelInfo()

	assert.NotEmpty(t, k.Kernel)
	assert.NotEmpty(t, k.Version)
	assert.NotEmpty(t, k.Hostname)
}
