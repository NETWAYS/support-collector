package icinga2

import (
	"fmt"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
	"github.com/NETWAYS/support-collector/internal/util"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/NETWAYS/support-collector/internal/collection"
)

const ModuleName = "icinga2"

var relevantPaths = []string{
	"/etc/icinga2",
	"/var/lib/icinga2",
}

var files = []string{
	"/etc/icinga2",
	"/var/lib/icinga2/api/packages/_api/active-stage",
}

var pluginFiles = []string{
	"/usr/lib64/nagios/plugins",
	"/usr/lib/nagios/plugins",
}

var optionalFiles = []string{
	"/etc/logrotate.d/icinga2",
	"/etc/icinga-installer/scenarios.d/last_scenario.yaml",
}

var detailedFiles = []string{
	"/var/log/icinga2/error.log",
	"/var/log/icinga2/crash",
	"/var/log/icinga2/debug.log",
	"/var/log/icinga2/icinga2.log",
	"/var/log/icinga-installer",
}

var commands = map[string][]string{
	"version.txt":                        {"icinga2", "-V"},
	"config-check.txt":                   {"icinga2", "daemon", "-C"},
	"objects-zones.txt":                  {"icinga2", "object", "list", "--type", "Zone"},
	"objects-endpoints.txt":              {"icinga2", "object", "list", "--type", "Endpoint"},
	"variables.txt":                      {"icinga2", "variable", "list"},
	"features.txt":                       {"icinga2", "feature", "list"},
	"user-icinga.txt":                    {"id", "icinga"},
	"user-nagios.txt":                    {"id", "nagios"},
	"icinga2-api-stage-directories.txt":  {"ls", "-ld", "/var/lib/icinga2/api/packages/_api/*/"},
	"director-api-stage-directories.txt": {"ls", "-ld", "/var/lib/icinga2/api/packages/director/*/"},
}

var detailedCommands = map[string][]string{
	"object-list.txt": {"icinga2", "object", "list"},
}

var possibleDaemons = []string{
	"/usr/lib/systemd/system/icinga2.service",
	"/etc/systemd/system/icinga2.service",
	"/etc/systemd/system/icinga2.service.d",
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewOutput(`(?i)(?:password|salt|token)\s*=\s*(.*)`, "icinga2", "variable"),
	obfuscate.NewFile(`(?i)(?:password|salt|token)\s*=\s*(.*)`, `conf`),
	obfuscate.NewFile(`(?i)(?:password|community)(.*)`, `log`),
}

func detectIcingaVersion(version string) string {
	result := regexp.MustCompile(`\(version:\s+r(\d+.\d+.\d+)`).FindStringSubmatch(version)

	return result[1]
}

func Collect(c *collection.Collection) {
	var icinga2version string

	if !util.ModuleExists(relevantPaths) {
		c.Log.Info("Could not find icinga2. Skipping")
		return
	}

	c.Log.Info("Collecting Icinga 2 information")

	c.RegisterObfuscators(obfuscators...)

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"),
		"*icinga2*",
		"netways-plugin*",
		"monitoring-plugin*",
		"nagios-*",
		"icinga-installer",
	)

	c.AddServiceStatusRaw(filepath.Join(ModuleName, "service.txt"), "icinga2")

	if collection.DetectServiceManager() == "systemd" {
		c.AddCommandOutput(filepath.Join(ModuleName, "systemd-icinga2.service"), "systemctl", "cat", "icinga2.service")
	}

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	c.AddFilesIfFound(ModuleName, pluginFiles...)

	for _, file := range optionalFiles {
		if _, err := os.Stat(file); err != nil {
			continue
		}

		c.AddFiles(ModuleName, file)
	}

	content, err := collection.LoadCommandOutput("icinga2", "-V")
	if err != nil {
		c.Log.Debug("Could not find executable for icinga2")

		icinga2version = ""
	} else {
		icinga2version = detectIcingaVersion(string(content))
	}

	// With Icinga 2 >= 2.14 the icinga2.debug cache is no longer built automatically on every reload. To retrieve a current state we build it manually (only possible from 2.14.0)
	// Needs to be done before commands are collected
	if icinga2version >= "2.14.0" {
		_, err = collection.LoadCommandOutput("icinga2", "daemon", "-C", "--dump-objects")
		if err != nil {
			c.Log.Warn(err)
		}
	}

	for name, cmd := range commands {
		c.AddCommandOutput(filepath.Join(ModuleName, name), cmd[0], cmd[1:]...)
	}

	for _, file := range possibleDaemons {
		c.AddFilesIfFound(ModuleName, file)
	}

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFilesIfFound(ModuleName, file)
		}

		for name, cmd := range detailedCommands {
			c.AddCommandOutput(filepath.Join(ModuleName, name), cmd[0], cmd[1:]...)
		}
	}

	// Collect from API endpoints if given
	if len(c.Config.Icinga2.Endpoints) > 0 {
		c.Log.Debug("Start to collect data from Icinga API endpoints")

		for _, e := range c.Config.Icinga2.Endpoints {
			c.Log.Debugf("New API endpoint found: '%s'. Trying...", e.Address)

			// Check if endpoint is reachable
			if err := e.IsReachable(5 * time.Second); err != nil { //nolint:mnd
				c.Log.Warn(err)
				continue
			}

			c.Log.Debug("Collect from resource 'v1/status'")

			// Request stats and health from endpoint
			res, err := e.Request("v1/status", 10*time.Second) //nolint:mnd
			if err != nil {
				c.Log.Warn(err)
				continue
			}

			// Save output to file. Replace "." in address with "_" and use as filename.
			c.AddFileJSON(filepath.Join(ModuleName, "api", "v1", "status", fmt.Sprintf("%s.json", strings.ReplaceAll(e.Address, ".", "_"))), res)

			c.Log.Debugf("Successfully finished endpoint '%s'", e.Address)
		}
	}
}
