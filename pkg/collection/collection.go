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

func (c *Collection) AddFiles(prefix, source string) {
	c.Log.Debug("Collecting files from ", source)

	files, err := LoadFiles(prefix, source)
	if err != nil {
		c.Log.Error(err)
	}

	c.Files = append(c.Files, files...)
}

func (c *Collection) AddCommandOutputWithTimeout(fileName string, timeout time.Duration, command string, arguments ...string) {
	c.Log.Debugf("Collecting command output: %s %v", command, arguments)

	output, err := LoadCommandOutputWithTimeout(timeout, command, arguments...)
	if err != nil {
		c.Log.Error(err)
	}

	c.AddFileData(fileName, output)

	return
}

func (c *Collection) AddCommandOutput(fileName, command string, arguments ...string) {
	c.AddCommandOutputWithTimeout(fileName, DefaultTimeout, command, arguments...)
}

func (c *Collection) AddInstalledPackagesRaw(fileName, pattern string) {
	c.Log.Debug("Collecting installed packages for pattern ", pattern)

	packages, err := ListInstalledPackagesRaw(pattern)
	if err != nil {
		c.Log.Warn(err)
	}

	c.AddFileData(fileName, packages)

	return
}

func (c *Collection) AddServiceStatusRaw(fileName, name string) {
	c.Log.Debug("Collecting service status for ", name)

	output, err := GetServiceStatusRaw(name)
	if err != nil {
		c.Log.Warn(err)
	}

	c.AddFileData(fileName, output)

	return
}
