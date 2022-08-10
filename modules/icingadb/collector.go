package icingadb

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"os"
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

var services = []string{
	"icingadb",
	"icingadb-redis",
	"icingadb-redis-server",
}

var journalctlLogs = map[string]collection.JournalElement{
	"journalctl-icingadb.txt":              {Service: "icingadb.service"},
	"journalctl-icingadb-redis.txt":        {Service: "icingadb-redis.service"},
	"journalctl-icingadb-redis-server.txt": {Service: "icingadb-redis-server.service"},
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

	c.AddInstalledPackagesRaw(ModuleName+"/packages.txt",
		"*icingadb*",
		"icingadb-redis",
	)

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for _, service := range services {
		c.AddServiceStatusRaw(ModuleName+"/service-"+service+".txt", service)
	}

	for name, element := range journalctlLogs {
		if service, err := collection.FindServices(element.Service); err == nil && len(service) > 0 {
			c.AddCommandOutput(ModuleName+"/"+name, "journalctl", "-u", element.Service, "--since", "7 days ago")
		}
	}
}
