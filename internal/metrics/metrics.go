package metrics

import (
	"github.com/NETWAYS/support-collector/internal/config"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
	"os"
	"strings"
	"time"
)

type Metrics struct {
	Command  string                   `json:"command"`
	Controls config.Config            `json:"controls"`
	Version  string                   `json:"version"`
	Timings  map[string]time.Duration `json:"timings"`
}

// New creates new Metrics
func New(version string) (m *Metrics) {
	return &Metrics{
		Command: getCommand(),
		Version: version,
		Timings: make(map[string]time.Duration),
	}
}

// getCommand returns the executed command and obfusactes *--icinga2* arguments
func getCommand() string {
	args := os.Args

	// Obfuscate icinga 2 api user and password
	for i, arg := range args {
		if strings.Contains(arg, "--icinga2") && i+1 < len(args) {
			args[i+1] = obfuscate.Replacement
		}
	}

	return strings.Join(args, " ")
}
