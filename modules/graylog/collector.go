package graylog

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"os"
	"path/filepath"
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
		c.Log.Info("Could not find Graylog")
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
