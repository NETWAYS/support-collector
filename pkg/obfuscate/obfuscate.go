package obfuscate

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Replacement is the standard replacement used during obfuscation.
const Replacement = "<HIDDEN>"

// Obfuscator is a generic interface a obfuscator should provide.
type Obfuscator interface {
	// IsAccepting should return true when the obfuscator should process the data.
	IsAccepting(...string) bool
	// Process takes an io.Reader, replaces all content to obfuscate and returns a new io.Reader.
	Process(io.Reader) (io.Reader, error)
}

func ProcessReader(r io.Reader, patterns []*regexp.Regexp) (out bytes.Buffer, err error) {
	rd := bufio.NewReader(r)

	var (
		line    string
		reading = true
	)

	for reading {
		line, err = rd.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
				reading = false

				if line == "" {
					break
				}
			} else {
				return out, fmt.Errorf("could not read from reader: %w", err)
			}
		}

		line = ReplacePatterns(line, patterns)

		_, _ = out.WriteString(line)
	}

	return out, err
}

func ReplacePatterns(line string, patterns []*regexp.Regexp) string {
	for _, pattern := range patterns {
		line = ReplacePattern(line, pattern)
	}

	return line
}

func ReplacePattern(line string, pattern *regexp.Regexp) string {
	return pattern.ReplaceAllStringFunc(line, func(s string) string {
		parts := pattern.FindStringSubmatch(s)

		if len(parts) > 1 {
			for _, match := range parts[1:] {
				if match != "" {
					s = strings.ReplaceAll(s, match, Replacement)
				}
			}

			return s
		}

		return Replacement
	})
}
