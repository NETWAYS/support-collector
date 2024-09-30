package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

const defaultAnswerFileName = "answer-file.yml"

// GenerateDefaultAnswerFile creates a new answer-file with default values in current dir
func GenerateDefaultAnswerFile() error {
	file, err := os.Create(defaultAnswerFileName)
	if err != nil {
		return fmt.Errorf("could not create answer file: %w", err)
	}

	defer file.Close()

	defaults := GetControlDefaultObject()

	yamlData, err := yaml.Marshal(defaults)
	if err != nil {
		return fmt.Errorf("could not marshal yamldata for answer-file: %w", err)
	}

	_, err = io.WriteString(file, string(yamlData))
	if err != nil {
		return fmt.Errorf("could not write to answer file: %w", err)
	}

	return nil
}

// ReadAnswerFile reads given values from answerFile and returns new Config
func ReadAnswerFile(answerFile string, conf *Config) error {
	data, err := os.ReadFile(answerFile)
	if err != nil {
		return fmt.Errorf("could not read answer file: %w", err)
	}

	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return fmt.Errorf("could not unmarshal answer file: %w", err)
	}

	return nil
}
