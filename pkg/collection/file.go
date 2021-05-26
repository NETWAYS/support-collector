package collection

import (
	"fmt"
	"io"
	"time"
)

type FileMap map[string]File

type File struct {
	Name     string
	Source   string
	Modified time.Time
	Data     []byte

	io.Writer
}

func NewFile(name string) *File {
	return &File{
		Name:     name,
		Modified: time.Now(),
	}
}

func NewFileFromFS(name, source string) (*File, error) {
	// TODO
	return nil, fmt.Errorf("not implemented")
}

func NewFileFromReader(name string, r io.Reader) (*File, error) {
	f := NewFile(name)

	_, err := io.Copy(f, r)
	if err != nil {
		err = fmt.Errorf("could not write to file buffer: %w", err)
	}

	return f, err
}

func (f *File) Write(p []byte) (n int, err error) {
	f.Data = append(f.Data, p...)

	return len(p), nil
}
