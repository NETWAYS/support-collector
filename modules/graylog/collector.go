package graylog

import (
	"github.com/NETWAYS/support-collector/internal/util"
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
)

const ModuleName = "graylog"

var relevantPaths = []string{
	"/etc/graylog",
}

var possibleDaemons = []string{
	"/usr/lib/systemd/system/graylog-server.service",
}

var files = []string{
	"/etc/graylog",
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:password_secret|root_password_sha2|elasticsearch_hosts|mongodb_uri).*\s*=\s*(.*)`, `conf`),
}

func Collect(c *collection.Collection) {
	if !util.ModuleExists(relevantPaths) {
		c.Log.Info("Could not find graylog. Skipping")
		return
	}

	c.Log.Info("Collecting Graylog information")

	c.RegisterObfuscators(obfuscators...)

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"), "*graylog*")
	c.AddServiceStatusRaw(filepath.Join(ModuleName, "service.txt"), "graylog-server")

	for _, file := range possibleDaemons {
		c.AddFilesIfFound(ModuleName, file)
	}

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}
}
