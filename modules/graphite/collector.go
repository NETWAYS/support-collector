package graphite

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"github.com/NETWAYS/support-collector/pkg/util"
	"os"
)

const ModuleName = "graphite"

var relevantPaths = []string{
	"/opt/graphite",
	"/etc/graphite-web",
}

var files = []string{
	"/opt/graphite/conf",
	"/opt/graphite/webapp/graphite/local_settings.py",
	"/etc/carbon",
	"/etc/graphite-api*",
	"/etc/graphite-web*",
}

var obfuscators = []*obfuscate.Obfuscator{
	// local_settings.py
	obfuscate.NewFile(`(?i)(?:USER|PASSWORD|HOST)\s*=\s*(.*)`, `py`),
	// *.conf
	obfuscate.NewFile(`(?i)(?:USER|PASSWORD|HOST)\s*=\s*(.*)`, `conf`),
}

var processFilter = []string{
	"graphite",
	"carbon",
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
		c.Log.Info("Could not find graphite")
		return
	}

	c.Log.Info("Collecting graphite information")

	c.RegisterObfuscators(obfuscators...)

	pythonFound := false
	versions := []string{"", "2", "3"}

	for _, suffix := range versions {
		if !util.CommandExists("python" + suffix) {
			continue
		}

		pythonFound = true

		commandsPython := map[string][]string{
			"python" + suffix + "-version.txt": {"python" + suffix, "--version"},
			"pip" + suffix + "-version.txt":    {"pip" + suffix, "--version"},
			"pip" + suffix + "-list.txt":       {"pip" + suffix, "list"},
		}

		for name, cmd := range commandsPython {
			c.AddCommandOutput(ModuleName+"/"+name, cmd[0], cmd[1:]...)
		}

		c.AddInstalledPackagesRaw(ModuleName+"/packages-python"+suffix+".txt", "*python"+suffix+"*")
	}

	if !pythonFound {
		c.Log.Warn("Python not found on system")
	}

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	processList, err := collection.ProcessListFilter(processFilter)
	if err != nil {
		c.Log.Warn("cant get process list")
	}

	// save process names to string array
	var processes string

	for _, process := range processList {
		processes = processes + process.Executable() + "\n"
	}

	c.AddFileDataRaw(ModuleName+"/processlist.txt", []byte(processes))

	c.AddInstalledPackagesRaw(ModuleName+"/packages-graphite.txt", "*graphite*")
}
