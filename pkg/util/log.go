package util

import (
	"github.com/sirupsen/logrus"
	"io"
)

type ExtraLogHook struct {
	Formatter logrus.Formatter
	Writer    io.Writer
}

func (h ExtraLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h ExtraLogHook) Fire(entry *logrus.Entry) error {
	line, err := h.Formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = h.Writer.Write(line)
	return err
}
