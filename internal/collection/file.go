package collection

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type File struct {
	Name     string
	Source   string
	Modified time.Time
	Data     []byte

	io.Writer
}

var reIgnoreFiles = regexp.MustCompile(`(^\.|~$|\.key$|\.csr$|\.crt$|\.pem$)`)

func NewFile(name string) *File {
	return &File{
		Name:     name,
		Modified: time.Now(),
	}
}

func LoadFiles(prefix, source string) (files []*File, err error) {
	// Is it a globbing pattern?
	if strings.ContainsAny(source, "*?") {
		return LoadFilesFromGlob(prefix, source)
	}

	stat, err := os.Stat(source)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("file does not exist '%s': %w", source, err)
			return
		}

		err = fmt.Errorf("could not stat file '%s': %w", source, err)

		return
	}

	if stat.IsDir() {
		return LoadFilesFromDirectory(prefix, source)
	}

	file, err := loadFile(prefix, source, stat)
	if err != nil {
		return
	}

	if !reIgnoreFiles.MatchString(file.Name) {
		files = append(files, file)
	}

	return
}

func loadFile(prefix, source string, stat os.FileInfo) (file *File, err error) {
	file = &File{
		Name:     path.Join(prefix, source),
		Modified: stat.ModTime(),
		Source:   source,
	}

	file.Data, err = os.ReadFile(source)
	if err != nil {
		err = fmt.Errorf("could not read file '%s': %w", source, err)
		return
	}

	return
}

func LoadFilesFromGlob(prefix, source string) (files []*File, err error) {
	var matches []string

	matches, err = filepath.Glob(source)
	if err != nil {
		err = fmt.Errorf("could not glob '%s': %w", source, err)
		return
	} else if len(matches) == 0 {
		err = fmt.Errorf("no files found for glob: '%s'", source) //nolint:goerr113
		return
	}

	for _, match := range matches {
		var globFiles []*File

		globFiles, err = LoadFiles(prefix, match)
		if err != nil {
			return
		}

		files = append(files, globFiles...)
	}

	return
}

//nolint:nakedret
func LoadFilesFromDirectory(prefix, source string) (files []*File, err error) {
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("could not walk path %s: %w", path, err)
		}

		// ignore directories or files matching the ignore pattern
		if info.IsDir() || reIgnoreFiles.MatchString(info.Name()) {
			return nil
		}

		var file *File

		// Document symlinks as text files
		if IsSymlink(info) {
			link, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("could not read link: %w", err)
			}

			file = NewFile(filepath.Join(prefix, path) + "-symlink.txt")
			file.Data = []byte(link)
		} else {
			file, err = loadFile(prefix, path, info)
			if err != nil {
				return err
			}
		}

		files = append(files, file)

		return nil
	})
	if err != nil {
		err = fmt.Errorf("error walking the path %s: %w", source, err)
		return
	}

	return
}

func IsSymlink(info os.FileInfo) bool {
	return info.Mode()&os.ModeSymlink == os.ModeSymlink
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
