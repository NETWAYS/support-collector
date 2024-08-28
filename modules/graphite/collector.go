package graphite

import (
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
	"github.com/NETWAYS/support-collector/internal/util"
)

const ModuleName = "graphite"

var relevantPaths = []string{
	"/opt/graphite",
	"/etc/graphite-web",
	"/etc/carbon",
}

var files = []string{
	"/opt/graphite/conf",
	"/opt/graphite/webapp/graphite/local_settings.py",
	"/etc/carbon",
	"/etc/graphite-api*",
	"/etc/graphite-web*",
	"/etc/sysconfig/graphite-api",
}

var detailedFiles = []string{
	"/var/log/carbon",
	"/var/log/graphite",
}

var journalctlLogs = map[string]collection.JournalElement{
	"journalctl-graphite-api.txt": {Service: "graphite-api.service"},
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

func Collect(c *collection.Collection) {
	if !util.ModuleExists(relevantPaths) {
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
			c.AddCommandOutput(filepath.Join(ModuleName, name), cmd[0], cmd[1:]...)
		}

		c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages-python"+suffix+".txt"), "*python"+suffix+"*")
	}

	if !pythonFound {
		c.Log.Debug("Python not found on system")
	}

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	processList, err := collection.ProcessListFilter(processFilter)
	if err != nil {
		c.Log.Debug("cant get process list")
	}

	// save process names to string array
	var processes string

	for _, process := range processList {
		processes = processes + process.Executable() + "\n"
	}

	c.AddFileDataRaw(filepath.Join(ModuleName, "processlist.txt"), []byte(processes))

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages-graphite.txt"), "*graphite*", "*carbon*")

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFilesIfFound(ModuleName, file)
		}

		for name, element := range journalctlLogs {
			if service, err := collection.FindServices(element.Service); err == nil && len(service) > 0 {
				c.AddJournalLog(filepath.Join(ModuleName, name), element.Service)
			}
		}
	}
}
