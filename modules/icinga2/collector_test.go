package icinga2

import (
	"bytes"
	"os"
	"testing"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
	"github.com/NETWAYS/support-collector/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	file, err := os.ReadFile("testdata/icinga-version.txt")
	if err != nil {
		t.Skip("cant read version file")
	}

	version := detectIcingaVersion(string(file))
	if version == "" {
		t.Skip("cant detect icinga2 version")
	}

	if !detectIcinga() {
		t.Skip("could not find icinga2 in the test environment")
		return
	}

	c := collection.New(&bytes.Buffer{})

	Collect(c)

	err = c.Close()
	assert.NoError(t, err)
}

func TestObfuscators(t *testing.T) {
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/etc/icinga2/constants.conf")
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/etc/icinga2/features-available/ido-mysql.conf")
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindFile, "/var/log/icinga2/debug.log")
	util.AssertObfuscationExample(t, obfuscators, obfuscate.KindOutput, "icinga2 variable list")

	util.AssertAllObfuscatorsTested(t, obfuscators)
}
