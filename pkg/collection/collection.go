package collection

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

type Collection struct {
	Files   []*File
	Log     *logrus.Logger
	LogData *bytes.Buffer
}

func New() (c *Collection) {
	c = &Collection{}
	c.LogData = &bytes.Buffer{}
	c.Log = logrus.New()
	c.Log.Out = c.LogData

	c.Log.Info("Initializing new collection")

	return
}

func (c *Collection) AddFileFromReader(name string, r io.Reader) (err error) {
	f, err := NewFileFromReader(name, r)
	if err != nil {
		return
	}

	c.Files = append(c.Files, f)

	return
}

func (c *Collection) AddFileData(fileName string, data []byte) {
	file := NewFile(fileName)
	file.Data = data

	c.Files = append(c.Files, file)
}

func (c *Collection) AddFiles(prefix, source string) (err error) {
	files, err := LoadFiles(prefix, source)
	if err != nil {
		return
	}

	c.Files = append(c.Files, files...)

	return
}

func (c *Collection) AddCommandOutputWithTimeout(fileName string, timeout time.Duration, command string, arguments ...string) (err error) {
	output, err := LoadCommandOutputWithTimeout(timeout, command, arguments...)
	// err is returned, but we add the file anyway

	c.AddFileData(fileName, output)

	return
}

func (c *Collection) AddCommandOutput(fileName, command string, arguments ...string) (err error) {
	return c.AddCommandOutputWithTimeout(fileName, DefaultTimeout, command, arguments...)
}
