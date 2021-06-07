package obfuscate

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestReplacePattern(t *testing.T) {
	assert.Equal(t, `password = <HIDDEN>`,
		ReplacePattern(`password = "XXX"`, regexp.MustCompile(`password\s*=\s*(.*)`)))

	assert.Equal(t, `<HIDDEN>`,
		ReplacePattern(`password = "XXX"`, regexp.MustCompile(`password\s*=.*`)))
}
