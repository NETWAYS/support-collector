package collection

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
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
	Obfuscators []*obfuscate.Obfuscator
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

// AddFileToOutput adds a single file while calling any obfuscator.
func (c *Collection) AddFileToOutput(file *File) (err error) {
	file.Data, err = c.callObfuscators(obfuscate.KindFile, file.Name, file.Data)
	if err != nil {
		c.Log.Warn(err)
	}

	err = c.AddFileToOutputRaw(file)
	if err != nil {
		c.Log.Warn(err)
	}

	return
}

// AddFileToOutputRaw adds the file directly to ZIP output.
//
// No obfuscation is applied.
func (c *Collection) AddFileToOutputRaw(file *File) (err error) {
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

func (c *Collection) AddFileFromReaderRaw(name string, r io.Reader) (err error) {
	f, err := NewFileFromReader(name, r)
	if err != nil {
		return
	}

	return c.AddFileToOutputRaw(f)
}

func (c *Collection) AddFileDataRaw(fileName string, data []byte) {
	file := NewFile(fileName)
	file.Data = data

	err := c.AddFileToOutputRaw(file)
	if err != nil {
		c.Log.Warn(err)
	}
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

	// obfuscate
	output, err = c.callObfuscators(obfuscate.KindOutput, obfuscate.JoinCommand(command, arguments...), output)
	if err != nil {
		c.Log.Warn(err)
	}

	c.AddFileDataRaw(file, output)
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

	c.AddFileDataRaw(fileName, packages)
}

func (c *Collection) AddServiceStatusRaw(fileName, name string) {
	c.Log.Debug("Collecting service status for ", name)

	output, err := GetServiceStatusRaw(name)
	if err != nil {
		c.Log.Warn(err)
	}

	c.AddFileDataRaw(fileName, output)
}

func (c *Collection) AddGitRepoInfo(fileName, path string) {
	c.Log.Debug("Collecting GIT repository details for ", path)

	info, err := LoadGitRepoInfo(path)
	if err != nil {
		c.Log.Warn(err)
	}

	c.AddFileYAML(fileName, info)
}

func (c *Collection) RegisterObfuscator(o *obfuscate.Obfuscator) {
	c.Obfuscators = append(c.Obfuscators, o)
}

func (c *Collection) ClearObfuscators() {
	c.Obfuscators = c.Obfuscators[:0]
}

// callObfuscators iterates all obfuscators and call them when matching.
func (c *Collection) callObfuscators(kind obfuscate.Kind, name string, data []byte) (out []byte, err error) {
	out = data

	for _, o := range c.Obfuscators {
		if o.IsAccepting(kind, name) {
			out, err = o.Process(data)
			if err != nil {
				return
			}
		}
	}

	return
}
