package postgresql

import (
	"bytes"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"testing"
)

func TestCollect(t *testing.T) {
	c := collection.New(&bytes.Buffer{})
	// c.Log = logrus.StandardLogger()

	Collect(c)
}
