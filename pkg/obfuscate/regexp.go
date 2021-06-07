package obfuscate

import "regexp"

// NewRegexpKeyValue builds an regexp.Regexp to match multiple key-value kinds.
//
// Should work with INI files and Icinga 2 config.
func NewRegexpKeyValue(key string) *regexp.Regexp {
	return regexp.MustCompile(`(?i)^.*` + key + `\s*=\s*(.*)$`)
}
