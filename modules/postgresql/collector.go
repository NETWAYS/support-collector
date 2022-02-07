package postgresql

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
)

const ModuleName = "postgresql"

var relevantPaths = []string{
	"/etc/postgresql",
	"/var/lib/pgsql",
}

var files = []string{
	"/etc/postgresql*",
	"/var/lib/pgsql*",
}

var possibleServices = []string{
	"postgresql",
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
		c.Log.Info("Could not find PostgreSQL")
		return
	}

	c.Log.Info("Collecting PostgreSQL information")

	c.AddInstalledPackagesRaw(ModuleName+"/packages.txt", "*postgresql*", "*pgsql*")
	c.AddFilesAtLeastOne(ModuleName, files...)

	for _, service := range possibleServices {
		c.AddServiceStatusRaw(ModuleName+"/service-"+service+".txt", service)
	}
}
