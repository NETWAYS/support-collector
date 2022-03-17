package icinga2

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"os"
	"os/exec"
)

const ModuleName = "icinga2"

var files = []string{
	"/etc/icinga2",
	"/var/log/icinga2/icinga2.log",
	"/usr/lib64/nagios/plugins",
}

var optionalFiles = []string{
	"/var/log/icinga2/error.log",
	"/var/log/icinga2/crash",
	"/var/log/icinga2/debug.log",
}

var commands = map[string][]string{
	"version.txt":           {"icinga2", "-V"},
	"config-check.txt":      {"icinga2", "daemon", "-C"},
	"objects-zones.txt":     {"icinga2", "object", "list", "--type", "Zone"},
	"objects-endpoints.txt": {"icinga2", "object", "list", "--type", "Endpoint"},
	"variables.txt":         {"icinga2", "variable", "list"},
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewOutput(`(?i)(?:password|salt|token)\s*=\s*(.*)`, "icinga2", "variable"),
	obfuscate.NewFile(`(?i)(?:password|salt|token)\s*=\s*(.*)`, `conf`),
	obfuscate.NewFile(`(?i)(?:password|community)(.*)`, `log`),
}

func DetectIcinga() bool {
	_, err := exec.LookPath("icinga2")
	return err == nil
}

func Collect(c *collection.Collection) {
	if !DetectIcinga() {
		c.Log.Info("Could not find icinga2")
		return
	}

	c.Log.Info("Collecting Icinga 2 information")

	c.RegisterObfuscators(obfuscators...)

	c.AddInstalledPackagesRaw(ModuleName+"/packages.txt",
		"*icinga2*",
		"netways-plugin*",
		"monitoring-plugin*",
		"nagios-plugin*")

	c.AddServiceStatusRaw(ModuleName+"/service.txt", "icinga2")

	if collection.DetectServiceManager() == "systemd" {
		c.AddCommandOutput(ModuleName+"/systemd-icinga2.service", "systemctl", "cat", "icinga2.service")
	}

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for _, file := range optionalFiles {
		if _, err := os.Stat(file); err != nil {
			continue
		}

		c.AddFiles(ModuleName, file)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(ModuleName+"/"+name, cmd[0], cmd[1:]...)
	}
}
