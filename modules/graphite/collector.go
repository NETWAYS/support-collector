package graphite

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/obfuscate"
	"github.com/NETWAYS/support-collector/pkg/util"
	"os"
	"path/filepath"
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
	"/var/log/carbon",
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
			c.AddCommandOutput(filepath.Join(ModuleName, name), cmd[0], cmd[1:]...)
		}

		c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages-python"+suffix+".txt"), "*python"+suffix+"*")
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

	timestamp := "7 days ago"

	for name, element := range journalctlLogs {
		if service, err := collection.FindServices(element.Service); err == nil && len(service) > 0 {
			c.AddCommandOutput(filepath.Join(ModuleName, name), "journalctl", "-u", element.Service, "--since", timestamp)
		}
	}

	c.AddFileDataRaw(filepath.Join(ModuleName, "processlist.txt"), []byte(processes))

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages-graphite.txt"), "*graphite*", "*carbon*")
}
