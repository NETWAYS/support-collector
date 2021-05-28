package collection

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

type Collection struct {
	Output  *zip.Writer
	Log     *logrus.Logger
	LogData *bytes.Buffer
}

func New(w io.Writer) (c *Collection) {
	c = &Collection{}
	c.LogData = &bytes.Buffer{}
	c.Log = logrus.New()
	c.Log.Out = c.LogData

	c.Output = zip.NewWriter(w)

	return
}

func (c *Collection) Close() error {
	return c.Output.Close()
}

func (c *Collection) AddFileToOutput(file *File) (err error) {
	fh := &zip.FileHeader{
		Name:     file.Name,
		Modified: file.Modified,
	}

	// Create file header
	fileWriter, err := c.Output.CreateHeader(fh)
	if err != nil {
		return fmt.Errorf("could not add file to zip: %w", err)
	}

	// Write data to ZIP
	_, err = io.Copy(fileWriter, bytes.NewReader(file.Data))
	if err != nil {
		return fmt.Errorf("could not write file to zip: %w", err)
	}

	return
}

func (c *Collection) AddLogToOutput() (err error) {
	if c.LogData == nil {
		return
	}

	fh := &zip.FileHeader{
		Name:     "support-collector.log",
		Modified: time.Now(),
	}
	logBuffer := bytes.NewBuffer(c.LogData.Bytes())

	if logBuffer.Len() != 0 {
		log, err := c.Output.CreateHeader(fh)
		if err != nil {
			return fmt.Errorf("could not add file to zip: %w", err)
		}

		_, err = io.Copy(log, logBuffer)
		if err != nil {
			return fmt.Errorf("could not write file to zip: %w", err)
		}
	}

	return
}

func (c *Collection) AddFileFromReader(name string, r io.Reader) (err error) {
	f, err := NewFileFromReader(name, r)
	if err != nil {
		return
	}

	return c.AddFileToOutput(f)
}

func (c *Collection) AddFileData(fileName string, data []byte) {
	file := NewFile(fileName)
	file.Data = data

	_ = c.AddFileToOutput(file)
}

func (c *Collection) AddFiles(prefix, source string) {
	c.Log.Debug("Collecting files from ", source)

	files, err := LoadFiles(prefix, source)
	if err != nil {
		c.Log.Error(err)
	}

	for _, file := range files {
		_ = c.AddFileToOutput(file)
	}
}

func (c *Collection) AddCommandOutputWithTimeout(file string,
	timeout time.Duration, command string, arguments ...string) {
	c.Log.Debugf("Collecting command output: %s %v", command, arguments)

	output, err := LoadCommandOutputWithTimeout(timeout, command, arguments...)
	if err != nil {
		c.Log.Error(err)
	}

	c.AddFileData(file, output)
}

func (c *Collection) AddCommandOutput(file, command string, arguments ...string) {
	c.AddCommandOutputWithTimeout(file, DefaultTimeout, command, arguments...)
}

func (c *Collection) AddInstalledPackagesRaw(fileName string, pattern ...string) {
	c.Log.Debug("Collecting installed packages for pattern ", strings.Join(pattern, " "))

	packages, err := ListInstalledPackagesRaw(pattern...)
	if err != nil {
		c.Log.Warn(err)
	}

	c.AddFileData(fileName, packages)
}

func (c *Collection) AddServiceStatusRaw(fileName, name string) {
	c.Log.Debug("Collecting service status for ", name)

	output, err := GetServiceStatusRaw(name)
	if err != nil {
		c.Log.Warn(err)
	}

	c.AddFileData(fileName, output)
}
