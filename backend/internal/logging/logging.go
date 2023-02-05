package logging

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type Hook struct {
	Formatter logrus.Formatter
	writer    *os.File
	levels    []logrus.Level
	logrus.Hook
}

func NewHook(outputName string, logLevels []logrus.Level, formatter logrus.Formatter) (*Hook, error) {
	var err error
	hook := Hook{
		levels:    logLevels,
		Formatter: formatter,
	}
	if outputName == "stdout" {
		hook.writer = os.Stdout
	} else if outputName == "stderr" {
		hook.writer = os.Stderr
	} else {
		hook.writer, err = os.OpenFile(outputName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			return nil, err
		}
	}
	return &hook, nil
}

func (h *Hook) Levels() []logrus.Level {
	return h.levels
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	line, err := h.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(line)
	if err != nil {
		return err
	}
	return nil
}

func GetLogger() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()
	}
	return logger
}

func InitLogger(configFile string) error {
	logger := GetLogger()
	cfg := LogConfig{}
	if configFile != "" {
		err := cfg.FromYAML(configFile)
		if err != nil {
			return err
		}
	}
	outputs, err := cfg.ListOutputsNames()
	if err != nil {
		return err
	}
	if cfg.DiscardDefaultLogging {
		logger.SetOutput(ioutil.Discard)
	}
	for _, output := range outputs {
		levels, err := cfg.GetLogLevels(output)
		if err != nil {
			return err
		}
		formatter, err := cfg.GetLogFormatter(output)
		if err != nil {
			return err
		}
		hook, err := NewHook(output, levels, formatter)
		if err != nil {
			return err
		}
		logger.AddHook(hook)
	}
	return nil
}
