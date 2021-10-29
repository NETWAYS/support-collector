package icinga2

import (
	"bytes"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"github.com/NETWAYS/support-collector/pkg/util"
	"testing"
)

func TestCollect(t *testing.T) {
	if !DetectIcinga() {
		t.Skip("could not find icinga2 in the test environment")
		return
	}

	c := collection.New(&bytes.Buffer{})
	// c.Log = logrus.StandardLogger()

	Collect(c)
}

func TestObfuscators(t *testing.T) {
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/etc/icinga2/constants.conf")
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/etc/icinga2/features-available/ido-mysql.conf")
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/var/log/icinga2/debug.log")
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindOutput, "icinga2 variable list")

	util.AssertAllObfuscatorsTested(t, obfuscators)
}
