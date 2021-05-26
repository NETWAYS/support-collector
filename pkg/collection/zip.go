package collection

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

func (c Collection) WriteZIP(w io.Writer) error {
	z := zip.NewWriter(w)

	for _, file := range c.Files {
		fh := &zip.FileHeader{
			Name:     file.Name,
			Modified: file.Modified,
		}

		// Create file header
		fileWriter, err := z.CreateHeader(fh)
		if err != nil {
			return fmt.Errorf("could not add file to zip: %w", err)
		}

		// Write data to ZIP
		_, err = io.Copy(fileWriter, bytes.NewReader(file.Data))
		if err != nil {
			return fmt.Errorf("could not write file to zip: %w", err)
		}
	}

	if len(c.Log) > 0 {
		log, err := z.Create("support-collector.log")
		if err != nil {
			return fmt.Errorf("could not add file to zip: %w", err)
		}

		_, err = io.Copy(log, bytes.NewReader(c.Log))
		if err != nil {
			return fmt.Errorf("could not write file to zip: %w", err)
		}
	}

	return z.Close()
}
