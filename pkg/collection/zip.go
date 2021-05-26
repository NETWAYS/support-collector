package collection

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"time"
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

	if c.LogData != nil {
		fh := &zip.FileHeader{
			Name:     "support-collector.log",
			Modified: time.Now(),
		}
		logBuffer := bytes.NewBuffer(c.LogData.Bytes())

		if logBuffer.Len() != 0 {
			log, err := z.CreateHeader(fh)
			if err != nil {
				return fmt.Errorf("could not add file to zip: %w", err)
			}

			_, err = io.Copy(log, logBuffer)
			if err != nil {
				return fmt.Errorf("could not write file to zip: %w", err)
			}
		}
	}

	return z.Close()
}
