package obfuscate

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

const iniExample = `[default]
user = "test"
password = "very-secret"
`

const iniResult = `[default]
user = "test"
password = <HIDDEN>
`

func TestFile_IsAccepting(t *testing.T) {
	o := NewFile(regexp.MustCompile(`^$`), ".ini")

	assert.True(t, o.IsAccepting("test.ini"))
	assert.False(t, o.IsAccepting("test.txt"))
}

func TestFile_Process(t *testing.T) {
	o := &File{
		Pattern: []*regexp.Regexp{
			regexp.MustCompile(`^password\s*=\s*(.*)`),
		},
	}

	out, err := o.Process(bytes.NewBufferString("default content\r\n"))
	assert.NoError(t, err)
	assert.Equal(t, "default content\r\n", out.String())

	out, err = o.Process(bytes.NewBufferString(iniExample))
	assert.NoError(t, err)
	assert.Equal(t, iniResult, out.String())
}
