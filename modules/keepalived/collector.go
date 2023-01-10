package keepalived

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"os"
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
		c.Log.Info("Could not find keepalived")
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
		c.AddServiceStatusRaw(ModuleName+"/service-"+service+".txt", service)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(ModuleName+"/"+name, cmd[0], cmd[1:]...)
	}
}
