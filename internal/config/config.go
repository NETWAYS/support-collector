package config

import (
	"fmt"
	"github.com/NETWAYS/support-collector/internal/util"
	"github.com/NETWAYS/support-collector/modules/icinga2/icingaapi"
	"slices"
	"time"
)

var (
	ModulesOrder = []string{
		"ansible",
		"base",
		"corosync",
		"elastic",
		"foreman",
		"grafana",
		"graphite",
		"graylog",
		"icinga-director",
		"icinga2",
		"icingadb",
		"icingaweb2",
		"influxdb",
		"keepalived",
		"mongodb",
		"mysql",
		"postgresql",
		"prometheus",
		"puppet",
		"redis",
		"webservers",
	}
)

type Config struct {
	General General `yaml:"general" json:"general"`
	Icinga2 Icinga2 `yaml:"icinga2" json:"icinga2"`
}

type General struct {
	OutputFile         string        `yaml:"outputFile" json:"outputFile"`
	AnswerFile         string        `yaml:"answerFile,omitempty" json:"answerFile,omitempty"`
	EnabledModules     []string      `yaml:"enabledModules" json:"enabledModules"`
	DisabledModules    []string      `yaml:"disabledModules" json:"disabledModules"`
	ExtraObfuscators   []string      `yaml:"extraObfuscators" json:"extraObfuscators"`
	DetailedCollection bool          `yaml:"detailedCollection" json:"detailedCollection"`
	CommandTimeout     time.Duration `yaml:"commandTimeout" json:"commandTimeout"`
}

type Icinga2 struct {
	Endpoints []icingaapi.Endpoint `yaml:"endpoints" json:"endpoints"`
}

// GetControlDefaultObject returns a new Config object with some pre-defined default values
func GetControlDefaultObject() Config {
	return Config{
		General: General{
			OutputFile:         util.BuildFileName(),
			AnswerFile:         "",
			EnabledModules:     []string{"all"},
			DisabledModules:    nil,
			ExtraObfuscators:   nil,
			DetailedCollection: true,
			CommandTimeout:     60 * time.Second,
		},
		Icinga2: Icinga2{},
	}
}

// ValidateConfig validates the given config.Config for errors. Returns []error
func ValidateConfig(conf Config) (errors []error) {
	for _, name := range conf.General.EnabledModules {
		if !slices.Contains(ModulesOrder, name) {
			errors = append(errors, fmt.Errorf("invalid module '%s' provided. Cant be enabled", name))
		}
	}
	return errors
}
