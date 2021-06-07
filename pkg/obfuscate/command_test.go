package obfuscate

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestCommand_IsAccepting(t *testing.T) {
	c := NewCommand(regexp.MustCompile(`^$`), "command", "sub-command")

	assert.True(t, c.IsAccepting("command", "sub-command"))
	assert.True(t, c.IsAccepting("command", "sub-command", "yolo"))
	assert.False(t, c.IsAccepting("command"))
	assert.False(t, c.IsAccepting("command", "other"))
}

func TestCommand_Process(t *testing.T) {
	c := NewCommand(regexp.MustCompile(`secret`), "command")

	out, err := c.Process(bytes.NewBufferString(`secret`))
	assert.NoError(t, err)
	assert.Equal(t, `<HIDDEN>`, out.String())
}
