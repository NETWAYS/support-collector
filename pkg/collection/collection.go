package collection

import (
	"io"
)

type Collection struct {
	Files []*File
	Log   []byte
}

func (c *Collection) AddFileFromReader(name string, r io.Reader) (err error) {
	f, err := NewFileFromReader(name, r)
	if err != nil {
		return
	}

	c.Files = append(c.Files, f)

	return
}

func (c *Collection) AddFiles(prefix, source string) (err error) {
	files, err := LoadFiles(prefix, source)
	if err != nil {
		return
	}

	c.Files = append(c.Files, files...)

	return
}
