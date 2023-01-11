package corosync

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
	"path/filepath"
)

const ModuleName = "corosync"

var relevantPaths = []string{
	"/etc/corosync",
}

var possibleDaemons = []string{
	"/lib/systemd/system/corosync.service",
	"/usr/lib/systemd/system/corosync.service",
	"/usr/lib/systemd/system/pacemaker.service",
}

var services = []string{
	"corosync",
	"pacemaker",
}

var files = []string{
	"/etc/corosync/corosync.conf",
	"/var/lib/pacemaker/cib/cib.xml",
	"/var/log/corosync/corosync.log",
	"/var/log/pacemaker/pacemaker.log",
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

	for _, service := range services {
		c.AddServiceStatusRaw(filepath.Join(ModuleName, "service-"+service+".txt"), service)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(filepath.Join(ModuleName, name), cmd[0], cmd[1:]...)
	}
}
