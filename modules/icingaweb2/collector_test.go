package icingaweb2

import (
	"bytes"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"github.com/NETWAYS/support-collector/pkg/util"
	"testing"
)

func TestCollect(t *testing.T) {
	if !Detect() {
		t.Skip("could not find icingaweb2 in the test environment")
		return
	}

	c := collection.New(&bytes.Buffer{})
	// c.Log = logrus.StandardLogger()

	Collect(c)
}

func TestObfuscators(t *testing.T) {
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/etc/icingaweb2/resources.ini")

	util.AssertAllObfuscatorsTested(t, obfuscators)
}
