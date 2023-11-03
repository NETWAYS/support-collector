package util

import (
	"github.com/sirupsen/logrus"
	"io"
)

type ExtraLogHook struct {
	Formatter logrus.Formatter
	Writer    io.Writer
	Level     logrus.Level
}

func (h ExtraLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h ExtraLogHook) Fire(entry *logrus.Entry) error {
	if entry.Level > h.Level {
		return nil
	}

	line, err := h.Formatter.Format(entry)
	if err != nil {
		return err
	}

	_, err = h.Writer.Write(line)

	return err
}
