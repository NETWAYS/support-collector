package prometheus

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"os"
	"path/filepath"
)

const ModuleName = "prometheus"

var relevantPaths = []string{
	"/etc/prometheus",
}

var possibleDaemons = []string{
	"/usr/lib/systemd/system/prometheus.service",
	"/lib/systemd/system/prometheus.service",
	"/usr/lib/systemd/system/pushgateway.service",
	"/lib/systemd/system/pushgateway.service",
	"/usr/lib/systemd/system/alertmanager.service",
	"/lib/systemd/system/alertmanager.service",
}

var services = []string{
	"alertmanager",
	"prometheus",
	"pushgateway",
}

var files = []string{
	"/etc/prometheus/prometheus.yml",
}

var commands = map[string][]string{
	"alertmanager-version.txt": {"alertmanager", "--version"},
	"prometheus-version.txt":   {"prometheus", "--version"},
	"pushgateway-version.txt":  {"pushgateway", "--version"},
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:password|secret).*\s*:\s*(.*)`, `yml`),
	obfuscate.NewFile(`(?i)(?:password|secret).*\s*:\s*(.*)`, `yaml`),
}

func Detect() bool {
	for _, path := range relevantPaths {
		_, err := os.Stat(path)
		if err == nil {
			return true
		}
	}

	return false
}

func Collect(c *collection.Collection) {
	if !Detect() {
		c.Log.Info("Could not find Prometheus")
		return
	}

	c.Log.Info("Collecting Prometheus information")

	c.RegisterObfuscators(obfuscators...)

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for _, file := range possibleDaemons {
		c.AddFilesIfFound(ModuleName, file)
	}

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"), "*prometheus*", "*pushgateway*", "*alertmanager*")

	for _, service := range services {
		c.AddServiceStatusRaw(filepath.Join(ModuleName, "service-"+service+".txt"), service)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(filepath.Join(ModuleName, name), cmd[0], cmd[1:]...)
	}
}
