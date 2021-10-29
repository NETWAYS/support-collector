package obfuscate

import (
	"regexp"
	"strings"
)

// NewRegexpKeyValue builds an regexp.Regexp to match multiple key-value kinds.
//
// Should work with INI files and Icinga 2 config.
func NewRegexpKeyValue(key string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)^.*` + key + `\s*=\s*(.*)$`)
}

func NewExtensionMatch(ext string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)\.` + ext + `$`)
}

func NewCommandMatch(command string, arguments ...string) *regexp.Regexp {
	return regexp.MustCompile(`^` + regexp.QuoteMeta(JoinCommand(command, arguments...)))
}

func JoinCommand(command string, arguments ...string) (s string) {
	s = command
	if len(arguments) > 0 {
		s += " " + strings.Join(arguments, " ")
	}

	return
}
