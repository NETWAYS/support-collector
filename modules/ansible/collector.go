package ansible

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
)

const ModuleName = "ansible"

var relevantPaths = []string{
	"/etc/ansible",
	"/usr/share/ansible",
}

var files = []string{
	"/etc/ansible",
}

var commands = map[string][]string{
	"version.txt": {"ansible", "--version"},
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
		c.Log.Info("Could not find ansible")
		return
	}

	c.Log.Info("Collecting Ansible information")

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(ModuleName+"/"+name, cmd[0], cmd[1:]...)
	}

	c.AddInstalledPackagesRaw(ModuleName+"/packages.txt", "*ansible*")
	c.AddInstalledPackagesRaw(ModuleName+"/packages-python.txt", "*python*")
	c.AddInstalledPackagesRaw(ModuleName+"/packages-pip.txt", "*pip*")
}
