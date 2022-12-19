package corosync

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
)

const ModuleName = "corosync"

var relevantPaths = []string{
	"/etc/corosync/service.d",
}

var possibleDaemons = []string{
	"/lib/systemd/system/corosync.service",
	"/usr/lib/systemd/system/corosync.service",
}

var files = []string{
	"/etc/corosync/corosync.conf",
	"/var/lib/pacemaker/cib/cib.xml",
}

var commands = map[string][]string{
	"version.txt": {"corosync", "-v"},
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
		c.Log.Info("Could not find corosync")
		return
	}

	c.Log.Info("Collecting corosync information")

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for _, file := range possibleDaemons {
		c.AddFilesIfFound(ModuleName, file)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(ModuleName+"/"+name, cmd[0], cmd[1:]...)
	}
}
