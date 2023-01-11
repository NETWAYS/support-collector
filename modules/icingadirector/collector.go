package icingadirector

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
	"path/filepath"
)

const (
	ModuleName       = "icinga_director"
	InstallationPath = "/usr/share/icingaweb2/modules/director"
)

var commands = map[string][]string{
	"health.txt":              {"icingacli", "director", "health"},
	"user-icingadirector.txt": {"id", "icingadirector"},
}

var possibleDaemons = []string{
	"/usr/lib/systemd/system/icinga-director.service",
	"/etc/systemd/system/icinga-director.service",
	"/etc/systemd/system/icinga-director.service.d",
}

var journalctlLogs = map[string]collection.JournalElement{
	"journalctl-director.txt": {Service: "icinga-director.service"},
}

// Detect if Icinga Director is installed on the system.
func Detect() bool {
	_, err := os.Stat(InstallationPath)
	return err == nil
}

// Collect data for Icinga Director.
func Collect(c *collection.Collection) {
	if !Detect() {
		c.Log.Info("Could not find Icinga Director")
		return
	}

	c.Log.Info("Collecting Icinga Director information")

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"), "*icinga*director*")
	c.AddServiceStatusRaw(filepath.Join(ModuleName, "service.txt"), "icinga-director")

	// TODO: more infos on modules, GIT details

	for name, cmd := range commands {
		c.AddCommandOutput(filepath.Join(ModuleName, name), cmd[0], cmd[1:]...)
	}

	for _, file := range possibleDaemons {
		c.AddFilesIfFound(ModuleName, file)
	}

	for name, element := range journalctlLogs {
		if service, err := collection.FindServices(element.Service); err == nil && len(service) > 0 {
			c.AddCommandOutput(filepath.Join(ModuleName, name), "journalctl", "-u", element.Service, "--since \"7 days ago\"")
		}
	}

	// Get GIT Repository details
	if path, ok := collection.IsGitRepository(InstallationPath); collection.DetectGitInstalled() && ok {
		c.AddGitRepoInfo(filepath.Join(ModuleName, "git-info.yml"), path)
	}
}
