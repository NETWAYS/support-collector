package keepalived

import (
	"github.com/NETWAYS/support-collector/internal/util"
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
)

const ModuleName = "keepalived"

var relevantPaths = []string{
	"/etc/keepalived",
}

var possibleDaemons = []string{
	"/usr/lib/systemd/system/keepalived.service",
}

var services = []string{
	"keepalived",
}

var files = []string{
	"/etc/keepalived/keepalived.conf",
}

var commands = map[string][]string{
	"version.txt": {"keepalived", "--version"},
}

var obfuscators = []*obfuscate.Obfuscator{
	// auth_pass in keepalived.conf
	obfuscate.NewFile(`(?i)(auth_pass) (.*)`, `conf`),
}

func Collect(c *collection.Collection) {
	if !util.ModuleExists(relevantPaths) {
		c.Log.Info("Could not find keepalived. Skipping")
		return
	}

	c.Log.Info("Collecting keepalived information")

	c.RegisterObfuscators(obfuscators...)

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
