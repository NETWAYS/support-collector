package puppet

import (
	"os"
	"path/filepath"

	"github.com/NETWAYS/support-collector/pkg/collection"
)

const ModuleName = "puppet"

var relevantPaths = []string{
	"/etc/puppetlabs",
	"/opt/puppetlabs",
}

var possibleServices = []string{
	"puppet",
	"puppetserver",
}

var detailedFiles = []string{
	"/var/log/puppet",
}

var files = []string{
	"/etc/puppetlabs",
	"/opt/puppetlabs/puppet/cache",
}

var commands = map[string][]string{
	"facter.txt":             {"facter"},
	"puppet-module-list.txt": {"puppet", "module", "list"},
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
		c.Log.Info("Could not find puppet")
		return
	}

	c.Log.Info("Collecting Puppet information")

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(filepath.Join(ModuleName, name), cmd[0], cmd[1:]...)
	}

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"), "*puppet*")

	for _, service := range possibleServices {
		c.AddServiceStatusRaw(filepath.Join(ModuleName, "service-"+service+".txt"), service)
	}

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFilesIfFound(ModuleName, file)
		}
	}
}
