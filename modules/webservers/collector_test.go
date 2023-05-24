package webservers

import (
	"bytes"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"testing"
)

func TestCollect(t *testing.T) {
	c := collection.New(&bytes.Buffer{})

	Collect(c)
}
