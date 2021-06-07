package obfuscate

import (
	"bytes"
	"io"
	"regexp"
	"strings"
)

// File provides a generic obfuscator for files with their name and contents.
type File struct {
	Obfuscator

	Extensions []string
	Pattern    []*regexp.Regexp
}

func NewFile(re *regexp.Regexp, extension string) *File {
	return &File{
		Extensions: []string{extension},
		Pattern:    []*regexp.Regexp{re},
	}
}

// IsAccepting checks if the file extension is matching a known list.
func (f File) IsAccepting(files ...string) bool {
	for _, file := range files {
		for _, ext := range f.Extensions {
			if strings.HasSuffix(file, ext) {
				return true
			}
		}
	}

	return false
}

// Process will replace all possible matches in an io.Reader.
func (f File) Process(r io.Reader) (bytes.Buffer, error) {
	return ProcessReader(r, f.Pattern)
}
