package icingaweb2

import (
	"github.com/NETWAYS/support-collector/internal/util"
	"os"
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
)

const (
	ModuleName  = "icingaweb2"
	ModulesPath = "/usr/share/icingaweb2/modules"
)

// Possible locations to indicate Icinga Web 2 is installed.
var relevantPaths = []string{
	"/etc/icingaweb2",
	"/usr/share/icingaweb2",
}

var files = []string{
	"/etc/icingaweb2",
}

var detailedFiles = []string{
	"/var/log/icingaweb2",
}

var journalctlLogs = map[string]collection.JournalElement{
	"journalctl-vspheredb.txt": {Service: "icinga-vspheredb.service"},
	"journalctl-reporting.txt": {Service: "icinga-reporting.service"},
	"journalctl-x509.txt":      {Service: "icinga-x509.service"},
}

var possibleDaemons = []string{
	"/usr/lib/systemd/system/icinga-vspheredb.service",
	"/etc/systemd/system/icinga-vspheredb.service",
	"/etc/systemd/system/icinga-vspheredb.service.d/",
	"/usr/lib/systemd/system/icinga-reporting.service",
	"/etc/systemd/system/icinga-reporting.service",
	"/etc/systemd/system/icinga-reporting.service.d",
	"/usr/lib/systemd/system/icinga-x509.service",
	"/etc/systemd/system/icinga-x509.service",
	"/etc/systemd/system/icinga-x509.service.d",
}

var tmpFiles = []string{
	"/usr/lib/tmpfiles.d/icinga-vspheredb.conf",
	"/etc/tmpfiles.d/icinga-vspheredb.conf",
}

var commands = map[string][]string{
	"version.txt":              {"icingacli", "version"},
	"modules.txt":              {"icingacli", "module", "list"},
	"vpsheredb-socket.txt":     {"ls", "-la", "/run/icinga-vspheredb/"},
	"user-icingavspheredb.txt": {"id", "icingavspheredb"},
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:bind_pw|password|token)\s*=\s*(.*)`, `ini`),
}

// Collect data for icingaweb2.
func Collect(c *collection.Collection) {
	if !util.ModuleExists(relevantPaths) {
		c.Log.Info("Could not find icingaweb2")
		return
	}

	c.Log.Info("Collecting Icinga Web 2 information")

	c.RegisterObfuscators(obfuscators...)

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"), "*icingaweb2*", "*icingacli*")

	if _, ok := collection.IsGitRepository("/usr/share/icingaweb2"); ok {
		c.AddGitRepoInfo(filepath.Join(ModuleName, "git.yml"), "/usr/share/icingaweb2")
	}

	CollectModuleInfo(c)

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(ModuleName+"/"+name, cmd[0], cmd[1:]...)
	}

	for _, file := range possibleDaemons {
		c.AddFilesIfFound(ModuleName, file)
	}

	for _, file := range tmpFiles {
		c.AddFilesIfFound(ModuleName, file)
	}

	// Detect PHP related packages and services
	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages-php.txt"), "*php*")

	if services, err := collection.FindServices("*php*-fpm"); err == nil && len(services) > 0 {
		for _, name := range services {
			c.AddServiceStatusRaw(filepath.Join(ModuleName, "service-"+name+".txt"), name)
		}
	}

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

	// Detect webserver packages
	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages-webserver.txt"), "*apache*", "*httpd*")
}

func CollectModuleInfo(c *collection.Collection) {
	if !collection.DetectGitInstalled() {
		c.Log.Debug("we need git to inspect modules closer")
	}

	modulesFiles, err := os.ReadDir(ModulesPath)
	if err != nil {
		c.Log.Debugf("Could not list modules in %s - %s", ModulesPath, err)
		return
	}

	for _, file := range modulesFiles {
		if !file.IsDir() {
			return
		}

		path := filepath.Join(ModulesPath, file.Name())

		if _, ok := collection.IsGitRepository(path); ok {
			c.AddGitRepoInfo(filepath.Join(ModuleName, "modules", file.Name()+".yml"), path)
		}
	}
}
