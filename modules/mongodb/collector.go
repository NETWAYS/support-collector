package mongodb

import (
	"os"
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
)

const ModuleName = "mongodb"

var relevantPaths = []string{
	"/etc/mongod.conf",
}

var possibleDaemons = []string{
	"/usr/lib/systemd/system/mongod.service",
	"/lib/systemd/system/mongod.service",
}

var services = []string{
	"mongod",
}

var files = []string{
	"/etc/mongod.conf",
}

var detailedFiles = []string{
	"/var/log/mongodb/mongod.log",
}

var commands = map[string][]string{
	"mongod-version.txt": {"mongod", "--version"},
	"mongo-version.txt":  {"mongo", "--version"},
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:password).*\s*:\s*(.*)`, `conf`),
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
		c.Log.Info("Could not find mongodb")
		return
	}

	c.Log.Info("Collecting mongodb information")

	c.RegisterObfuscators(obfuscators...)

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for _, file := range possibleDaemons {
		c.AddFilesIfFound(ModuleName, file)
	}

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"), "*mongo*")

	for _, service := range services {
		c.AddServiceStatusRaw(filepath.Join(ModuleName, "service-"+service+".txt"), service)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(filepath.Join(ModuleName, name), cmd[0], cmd[1:]...)
	}

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFilesIfFound(ModuleName, file)
		}
	}
}
