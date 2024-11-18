package redis

import (
	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
	"github.com/NETWAYS/support-collector/internal/util"
	"path/filepath"
)

const ModuleName = "redis"

var relevantPaths = []string{
	"/etc/redis",
	"/etc/redis.conf",
}

var files = []string{
	"/etc/redis*",
	"/proc/sys/vm/overcommit_memory",
}

var optionalFiles = []string{
	"/etc/logrotate.d/redis*",
}

var detailedFiles = []string{
	"/var/log/redis/redis-server.log",
	"/var/log/redis/redis.log",
}

var possibleDaemons = []string{
	"/lib/systemd/system/redis-server.service",
	"/lib/systemd/system/redis-server@.service",
	"/etc/systemd/system/redis*",
	"/usr/lib/systemd/system/redis*",
}

var services = []string{
	"redis-server",
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:requirepass)\s*(.*)`, `conf`),
}

func Collect(c *collection.Collection) {
	if !util.ModuleExists(relevantPaths) {
		c.Log.Info("Could not find redis. Skipping")
		return
	}

	c.Log.Info("Collecting redis information")

	c.RegisterObfuscators(obfuscators...)

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for _, file := range optionalFiles {
		c.AddFilesIfFound(ModuleName, file)
	}

	for _, file := range possibleDaemons {
		c.AddFilesIfFound(ModuleName, file)
	}

	for _, service := range services {
		c.AddServiceStatusRaw(filepath.Join(ModuleName, "service-"+service+".txt"), service)
	}

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFilesIfFound(ModuleName, file)
		}
	}
}
