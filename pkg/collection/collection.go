package collection

import (
	"io"
)

type Collection struct {
	Files []*File
	Log   []byte
}

func (c *Collection) AddFile(name string, r io.Reader) (err error) {
	f, err := NewFileFromReader(name, r)
	if err != nil {
		return
	}

	c.Files = append(c.Files, f)

	return
}
