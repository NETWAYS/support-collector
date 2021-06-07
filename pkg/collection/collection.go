package collection

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
	"time"
)

type Collection struct {
	Output      *zip.Writer
	Log         *logrus.Logger
	LogData     *bytes.Buffer
	ExecTimeout time.Duration
}

func New(w io.Writer) (c *Collection) {
	c = &Collection{
		Output:      zip.NewWriter(w),
		Log:         logrus.New(),
		LogData:     &bytes.Buffer{},
		ExecTimeout: DefaultTimeout,
	}

	c.Log.Out = c.LogData

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

func (c *Collection) AddFileYAML(fileName string, data interface{}) {
	var buf bytes.Buffer

	err := yaml.NewEncoder(&buf).Encode(&data)
	if err != nil {
		c.Log.Warnf("could not encode YAML data for '%s': %s", fileName, err)
	}

	file := NewFile(fileName)
	file.Data = buf.Bytes()

	_ = c.AddFileToOutput(file)
}

func (c *Collection) AddFiles(prefix, source string) {
	c.Log.Debug("Collecting files from ", source)

	files, err := LoadFiles(prefix, source)
	if err != nil {
		c.Log.Warn(err)
	}

	for _, file := range files {
		_ = c.AddFileToOutput(file)
	}
}

func (c *Collection) AddFilesAtLeastOne(prefix string, sources ...string) {
	var foundFiles int

	for _, source := range sources {
		files, _ := LoadFiles(prefix, source)
		if len(files) == 0 {
			return
		}

		c.Log.Debug("Collecting files from ", source)

		for _, file := range files {
			foundFiles++

			_ = c.AddFileToOutput(file)
		}
	}

	if foundFiles == 0 {
		c.Log.Warnf("Found no files under: %s", strings.Join(sources, " "))
	}
}

func (c *Collection) AddCommandOutputWithTimeout(file string,
	timeout time.Duration, command string, arguments ...string) {
	c.Log.Debugf("Collecting command output: %s %s", command, strings.Join(arguments, " "))

	output, err := LoadCommandOutputWithTimeout(timeout, command, arguments...)
	if err != nil {
		c.Log.Warn(err)
	}

	c.AddFileData(file, output)
}

func (c *Collection) AddCommandOutput(file, command string, arguments ...string) {
	c.AddCommandOutputWithTimeout(file, c.ExecTimeout, command, arguments...)
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

func (c *Collection) AddGitRepoInfo(fileName, path string) {
	c.Log.Debug("Collecting GIT repository details for ", path)

	info, err := LoadGitRepoInfo(path)
	if err != nil {
		c.Log.Warn(err)
	}

	c.AddFileYAML(fileName, info)
}
