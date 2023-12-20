package obfuscate

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Replacement is the standard replacement used during obfuscation.
const Replacement = "<HIDDEN>"

// Kind is used by Obfuscator to identify the kind of content to obfuscate.
type Kind uint8

const (
	// KindAny allows to obfuscate content in any resource.
	KindAny Kind = iota
	// KindFile declares a file resource.
	KindFile
	// KindOutput declares an output resource, e.g. from an command.
	KindOutput
)

// Obfuscator provides the basic functionality of an obfuscation engine.
//
// Kind filters the variant of resource we want to work on, while Affecting defines a list of regexp.Regexp, to match
// against for the file names, or command.
//
// Replacements will be iterated, so all matches or matched groups will be replaced.
type Obfuscator struct {
	Kind
	Affecting    []*regexp.Regexp
	Replacements []*regexp.Regexp
	Files        uint
	Replaced     uint
}

// New returns a basic Obfuscator with provided regexp.Regexp.
func New(kind Kind, affects, replace *regexp.Regexp) *Obfuscator {
	return &Obfuscator{
		Kind:         kind,
		Affecting:    []*regexp.Regexp{affects},
		Replacements: []*regexp.Regexp{replace},
	}
}

// NewFile returns an Obfuscator and will initialize regexp.Regexp based on extension and a string for replacement.
func NewFile(replace, ext string) *Obfuscator {
	o := &Obfuscator{Kind: KindFile}
	o.WithAffecting(NewExtensionMatch(ext))
	o.WithReplacement(regexp.MustCompile(replace))

	return o
}

// NewOutput returns an Obfuscator and will initialize regexp.Regexp based on command and replacement.
func NewOutput(replace, command string, arguments ...string) *Obfuscator {
	o := &Obfuscator{Kind: KindOutput}
	o.WithAffecting(NewCommandMatch(command, arguments...))
	o.WithReplacement(regexp.MustCompile(replace))

	return o
}

// NewAny returns an Obfuscator that can be used to replace any input.
func NewAny(replace string) *Obfuscator {
	o := &Obfuscator{Kind: KindAny}
	o.WithAffecting(regexp.MustCompile(`(.*)`))
	o.WithReplacement(regexp.MustCompile(replace))

	return o
}

// WithAffecting adds a new element to the list.
func (o *Obfuscator) WithAffecting(a *regexp.Regexp) *Obfuscator {
	o.Affecting = append(o.Affecting, a)

	return o
}

// WithReplacement adds a new element to the list.
func (o *Obfuscator) WithReplacement(r *regexp.Regexp) *Obfuscator {
	o.Replacements = append(o.Replacements, r)

	return o
}

// IsAccepting checks if we want to work on the resource.
func (o Obfuscator) IsAccepting(kind Kind, name string) bool {
	if o.Kind != KindAny && o.Kind != kind {
		return false
	}

	for _, re := range o.Affecting {
		if re.MatchString(name) {
			return true
		}
	}

	return false
}

// Process takes data and returns it obfuscated.
func (o *Obfuscator) Process(data []byte) (uint, []byte, error) {
	count, out, err := o.ProcessReader(bytes.NewReader(data))

	//goland:noinspection GoNilness
	return count, out.Bytes(), err
}

// ProcessReader takes an io.Reader and returns a new one obfuscated.
func (o *Obfuscator) ProcessReader(r io.Reader) (count uint, out bytes.Buffer, err error) {
	rd := bufio.NewReader(r)

	var (
		line    string
		ending  string
		reading = true
		c       uint
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
				return count, out, fmt.Errorf("could not read from reader: %w", err)
			}
		}

		// Pop any line ending before replacing
		ending = ""
		line = strings.TrimRightFunc(line, func(r rune) bool {
			if r == '\r' || r == '\n' {
				ending = string(r) + ending
				return true
			}
			return false
		})

		// Replace any matches, but skip empty lines
		if strings.TrimSpace(line) != "" {
			line, c = ReplacePatterns(line, o.Replacements)
		}

		// Add line ending back after replacement
		line += ending

		if c > 0 {
			count += c
			o.Replaced += c
		}

		_, _ = out.WriteString(line)
	}

	if count > 0 {
		o.Files++
	}

	return count, out, err
}

// ReplacePatterns replaces all the patterns matches in a line.
func ReplacePatterns(line string, patterns []*regexp.Regexp) (s string, count uint) {
	for _, pattern := range patterns {
		var c uint

		line, c = ReplacePattern(line, pattern)
		count += c
	}

	return line, count
}

// ReplacePattern replaces all matches in a line.
func ReplacePattern(line string, pattern *regexp.Regexp) (s string, count uint) {
	return pattern.ReplaceAllStringFunc(line, func(s string) string {
		parts := pattern.FindStringSubmatch(s)

		if len(parts) > 1 {
			for _, match := range parts[1:] {
				if match != "" {
					count++

					s = strings.ReplaceAll(s, match, Replacement)

					log.Debugf("replaced token for %s", pattern.String())
				}
			}

			return s
		}

		count++

		log.Debugf("replaced token for %s", pattern.String())

		return Replacement
	}), count
}
