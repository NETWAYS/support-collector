package base

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollect(t *testing.T) {
	c := collection.New()

	Collect(c)

	assert.NotEmpty(t, c.Files)
}