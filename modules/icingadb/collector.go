package icingadb

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"os"
	"path/filepath"
)

const (
	ModuleName = "icingadb"
)

var relevantPaths = []string{
	"/etc/icingadb",
}

var files = []string{
	"/etc/icingadb",
	"/etc/icingadb-redis",
	"/etc/icinga2/features-enabled/icingadb.conf",
}

var detailedFiles = []string{
	"/var/log/icingadb/",
	"/var/log/icingadb-redis",
}

var journalctlLogs = map[string]collection.JournalElement{
	"journalctl-icingadb.txt":              {Service: "icingadb.service"},
	"journalctl-icingadb-redis.txt":        {Service: "icingadb-redis.service"},
	"journalctl-icingadb-redis-server.txt": {Service: "icingadb-redis-server.service"},
}

var optionalFiles = []string{
	"/etc/logrotate.d/icingadb-redis-server",
}

var services = []string{
	"icingadb",
	"icingadb-redis",
	"icingadb-redis-server",
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:password)\s*=\s*(.*)`, `conf`),
	obfuscate.NewFile(`(?i)(?:password)\s*=\s*(.*)`, `yml`),
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
		c.Log.Info("Could not find IcingaDB")
		return
	}

	c.Log.Info("Collecting IcingaDB information")

	c.RegisterObfuscators(obfuscators...)

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"),
		"*icingadb*",
		"icingadb-redis",
	)

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for _, file := range optionalFiles {
		if _, err := os.Stat(file); err != nil {
			continue
		}

		c.AddFiles(ModuleName, file)
	}

	for _, service := range services {
		c.AddServiceStatusRaw(filepath.Join(ModuleName, "service-"+service+".txt"), service)
	}

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFilesIfFound(ModuleName, file)
		}

		for name, element := range journalctlLogs {
			if service, err := collection.FindServices(element.Service); err == nil && len(service) > 0 {
				c.AddCommandOutput(filepath.Join(ModuleName, name), "journalctl", "-u", element.Service, "--since", "7 days ago")
			}
		}
	}
}
