package util

import (
	"testing"

	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"github.com/stretchr/testify/assert"
)

var (
	noObfuscator []*obfuscate.Obfuscator
)

func testObfuscator() []*obfuscate.Obfuscator {
	return []*obfuscate.Obfuscator{
		obfuscate.NewOutput(`.*`, "command"),
		obfuscate.NewFile(`.*`, "log"),
	}
}

func TestLoadTestdata(t *testing.T) {
	mockT := new(testing.T)

	assert.Equal(t, "output\n", LoadTestdata(mockT, "command.txt"))
	assert.False(t, mockT.Failed())

	LoadTestdata(mockT, "nonexisting.txt")
	assert.True(t, mockT.Failed())
}

func TestAssertObfuscation(t *testing.T) {
	mockT := new(testing.T)
	AssertObfuscation(mockT, noObfuscator, obfuscate.KindFile, "a", "b", "c")
	assert.True(t, mockT.Failed())

	o := testObfuscator()

	mockT = new(testing.T)
	AssertObfuscation(mockT, o, obfuscate.KindOutput, "command", "b", "<HIDDEN>")
	AssertObfuscation(mockT, o, obfuscate.KindFile, "test.log", "b", "<HIDDEN>")
	assert.False(t, mockT.Failed())
}

func TestAssertObfuscationExample(t *testing.T) {
	mockT := new(testing.T)
	o := testObfuscator()

	AssertObfuscationExample(mockT, o, obfuscate.KindOutput, "command")
	AssertObfuscationExample(t, o, obfuscate.KindFile, "file/test.log")
	assert.False(t, mockT.Failed())
}

func TestAssertAllObfuscatorsTested(t *testing.T) {
	mockT := new(testing.T)

	AssertAllObfuscatorsTested(mockT, noObfuscator)
	assert.True(t, mockT.Failed())

	mockT = new(testing.T)
	o := testObfuscator()

	AssertObfuscation(mockT, o, obfuscate.KindOutput, "command", "b", "<HIDDEN>")
	AssertObfuscation(mockT, o, obfuscate.KindFile, "test.log", "b", "<HIDDEN>")
	AssertAllObfuscatorsTested(t, o)
	assert.False(t, mockT.Failed())

	mockT = new(testing.T)

	AssertAllObfuscatorsTested(mockT, testObfuscator())
	assert.True(t, mockT.Failed())
}
