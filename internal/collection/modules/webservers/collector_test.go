package webservers

import (
	"bytes"
	"github.com/NETWAYS/support-collector/internal/collection"
	"testing"
)

func TestCollect(t *testing.T) {
	c := collection.New(&bytes.Buffer{})

	Collect(c)
}
