package obfuscate

import (
	"bytes"
	"io"
	"regexp"
)

// Command provides a generic obfuscator for command line, matching on command and replacing output.
type Command struct {
	Obfuscator

	Command []string
	Pattern []*regexp.Regexp
}

func NewCommand(re *regexp.Regexp, command ...string) *Command {
	return &Command{
		Command: command,
		Pattern: []*regexp.Regexp{re},
	}
}

// IsAccepting checks if the command line is at least partially matching.
func (f Command) IsAccepting(command ...string) bool {
	if len(f.Command) > len(command) {
		return false
	}

	for i, v := range f.Command {
		if v != command[i] {
			return false
		}
	}

	return true
}

// Process will replace all possible matches in an io.Reader.
func (f Command) Process(r io.Reader) (bytes.Buffer, error) {
	return ProcessReader(r, f.Pattern)
}
