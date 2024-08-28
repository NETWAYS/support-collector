package foreman

import (
	"github.com/NETWAYS/support-collector/internal/util"
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
)

const ModuleName = "foreman"

var relevantPaths = []string{
	"/etc/foreman",
}

var files = []string{
	"/etc/foreman",
	"/etc/foreman-installer",
	"/etc/foreman-proxy",
}

var detailedFiles = []string{
	"/var/log/foreman",
	"/var/log/foreman-installer",
	"/var/log/foreman-proxy",
}

var obfuscaters = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:password)\s*:\s*(.*)`, "yml"),
	obfuscate.NewFile(`(?i)(?:ENCRYPTION_KEY)\s*=\s*(.*)`, "rb"),
}

func Collect(c *collection.Collection) {
	if !util.ModuleExists(relevantPaths) {
		c.Log.Info("Could not find Foreman")
		return
	}

	c.Log.Info("Collection Foreman information")

	c.RegisterObfuscators(obfuscaters...)

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"),
		"foreman",
		"foreman-installer",
		"foreman-proxy",
	)

	c.AddServiceStatusRaw(filepath.Join(ModuleName, "service.txt"), "foreman")

	if collection.DetectServiceManager() == "systemd" {
		c.AddCommandOutput(filepath.Join(ModuleName, "systemd-foreman.service"), "systemctl", "cat", "foreman.service")
	}

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFiles(ModuleName, file)
		}
	}
}
