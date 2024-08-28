package postgresql

import (
	"github.com/NETWAYS/support-collector/internal/util"
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
)

const ModuleName = "postgresql"

var relevantPaths = []string{
	"/etc/postgresql",
	"/var/lib/pgsql",
}

var files = []string{
	"/etc/postgresql*",
	"/var/lib/pgsql/data/*.conf", // RedHat based systems, where the configuration is found
}

var detailedFiles = []string{
	"/var/log/postgresql",
	"/var/lib/pgsql/pgstartup.log",
}

var commands = map[string][]string{
	"version.txt": {"psql", "-V"},
}

var possibleServices = []string{
	"postgresql",
}

func Collect(c *collection.Collection) {
	if !util.ModuleExists(relevantPaths) {
		c.Log.Info("Could not find PostgreSQL")
		return
	}

	c.Log.Info("Collecting PostgreSQL information")

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"), "*postgresql*", "*pgsql*")
	c.AddFilesIfFound(ModuleName, files...)

	for _, service := range possibleServices {
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
