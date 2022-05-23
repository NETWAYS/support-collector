package puppet

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
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

var files = []string{
	"/etc/puppetlabs",
	"/opt/puppetlabs/puppet/cache",
}

var commands = map[string][]string{
	"version-puppet.txt":       {"puppet", "--version"},
	"version-puppetserver.txt": {"puppetserver", "--version"},
	"puppet-module-list.txt":   {"puppet", "module", "list"},
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
		c.AddCommandOutput(ModuleName+"/"+name, cmd[0], cmd[1:]...)
	}

	c.AddInstalledPackagesRaw(ModuleName+"/packages.txt", "*puppet*")

	for _, service := range possibleServices {
		c.AddServiceStatusRaw(ModuleName+"/service-"+service+".txt", service)
	}
}
