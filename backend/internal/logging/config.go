package logging

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type LogConfig struct {
	Outputs               []Output `yaml:"outputs"`
	DiscardDefaultLogging bool     `yaml:"discard_default_logging"`
}

type Output struct {
	OutputName string   `yaml:"output_name"`
	Levels     []string `yaml:"levels"`
	Formatter  string   `yaml:"formatter"`
}

func (c *LogConfig) FromYAML(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return err
	}

	return nil
}

func (c *LogConfig) GetLogLevels(outputName string) ([]logrus.Level, error) {
	for _, output := range c.Outputs {
		if output.OutputName == outputName {
			levels := []logrus.Level{}
			for _, textLevel := range output.Levels {
				level, err := logrus.ParseLevel(textLevel)
				if err != nil {
					return nil, err
				}
				levels = append(levels, level)
			}
			return levels, nil
		}
	}
	return nil, fmt.Errorf("Hook with output name %s does not exist", outputName)
}

func (c *LogConfig) ListOutputsNames() ([]string, error) {
	outputs := []string{}
	for _, output := range c.Outputs {
		outputs = append(outputs, output.OutputName)
	}
	return outputs, nil
}

func (c *LogConfig) GetLogFormatter(outputName string) (logrus.Formatter, error) {
	for _, output := range c.Outputs {
		if output.OutputName == outputName {
			switch strings.ToLower(output.Formatter) {
			case "json":
				return &logrus.JSONFormatter{}, nil

			case "text":
				return &logrus.TextFormatter{}, nil

			default:
				return nil, fmt.Errorf("%s is not a valid formatter name for hook %s", output.Formatter, outputName)
			}
		}
	}
	return nil, fmt.Errorf("Hook with output name %s does not exist", outputName)
}
