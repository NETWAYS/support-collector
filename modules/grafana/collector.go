package grafana

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"os"
)

const ModuleName = "grafana"

var relevantPaths = []string{
	"/etc/grafana",
	"/usr/share/grafana",
}

var files = []string{
	"/etc/grafana",
}

var commands = map[string][]string{
	"grafana-cli-version.txt":      {"grafana-cli", "-v"},
	"grafana-cli-plugins-list.txt": {"grafana-cli", "plugins", "ls"},
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:password|token)\s*=\s*(.*)`, `ini`),
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
		c.Log.Info("Could not find grafana")
		return
	}

	c.Log.Info("Collecting grafana information")

	c.RegisterObfuscators(obfuscators...)

	c.AddInstalledPackagesRaw(ModuleName+"/packages.txt", "*grafana*", "*chrome*", "*chromium*")
	c.AddServiceStatusRaw(ModuleName+"/service.txt", "grafana-server")

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(ModuleName+"/"+name, cmd[0], cmd[1:]...)
	}
}
