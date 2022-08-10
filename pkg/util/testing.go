package util

import (
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

// AssertObfuscation is a helper function for tests where we want to validate if obfuscation.Obfuscator works correctly.
func AssertObfuscation(t *testing.T, obfuscators []*obfuscate.Obfuscator,
	kind obfuscate.Kind, name, input, expected string) {
	t.Helper()

	for _, o := range obfuscators {
		if !o.IsAccepting(kind, name) {
			continue
		}

		_, out, err := o.Process([]byte(input))
		if err != nil {
			t.Errorf("error during obfuscation: %s", err)
			return
		}

		assert.Equal(t, expected, string(out))

		return
	}

	t.Errorf("no obfuscator found for: %s", name)
}

// AssertObfuscationExample uses AssertObfuscation to assert but loads the example automatically from testdata.
//
// Parameter `name` must correspond to the relative file name under testdata, or in case of a command, spaces are
// replaced by a minus sign, and txt is used for the file extension.
func AssertObfuscationExample(t *testing.T, obfuscators []*obfuscate.Obfuscator, kind obfuscate.Kind, name string) {
	t.Helper()

	var path string

	switch kind { //nolint:exhaustive
	case obfuscate.KindFile:
		path = strings.TrimPrefix(name, "/")
	case obfuscate.KindOutput:
		path = strings.ReplaceAll(name, " ", "-") + ".txt"
	default:
		t.Errorf("AssertObfuscationExample not implemented for kind %T", kind)
		return
	}

	AssertObfuscation(t, obfuscators, kind, name, LoadTestdata(t, path), LoadTestdata(t, path+".obfuscated"))
}

func AssertAllObfuscatorsTested(t *testing.T, obfuscators []*obfuscate.Obfuscator) {
	t.Helper()

	var all, missing uint

	for _, o := range obfuscators {
		all++

		if o.Replaced == 0 {
			missing++
		}
	}

	if missing > 0 {
		t.Errorf("%d of %d obfuscators where not triggered", missing, all)
		return
	} else if all == 0 {
		t.Error("no obfuscator defined")
		return
	}
}

// LoadTestdata loads a file from the testdata directory and returns its contents as string.
//
// Intended to load text file for comparison in assertions.
func LoadTestdata(t *testing.T, name string) string {
	t.Helper()

	content, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Errorf("could not load testdata file: %s - %s", name, err)
		return ""
	}

	return string(content)
}
