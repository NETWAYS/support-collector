package icingadirector

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
)

const (
	ModuleName       = "icinga_director"
	InstallationPath = "/usr/share/icingaweb2/modules/director"
)

var commands = map[string][]string{
	"health.txt": {"icingacli", "director", "health"},
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

	c.AddInstalledPackagesRaw(ModuleName+"/packages.txt", "*icinga*director*")
	c.AddServiceStatusRaw(ModuleName+"/service.txt", "icinga-director")

	// TODO: more infos on modules, GIT details

	for name, cmd := range commands {
		c.AddCommandOutput(ModuleName+"/"+name, cmd[0], cmd[1:]...)
	}

	// Get GIT Repository details
	if path, ok := collection.IsGitRepository(InstallationPath); collection.DetectGitInstalled() && ok {
		c.AddGitRepoInfo(ModuleName+"/git-info.yml", path)
	}
}
