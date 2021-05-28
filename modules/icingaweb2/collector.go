package icingaweb2

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"io/ioutil"
	"os"
	"path/filepath"
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
	"/var/log/icingaweb2",
}

var commands = map[string][]string{
	"version.txt": {"icingacli", "version"},
	"modules.txt": {"icingacli", "module", "list"},
}

// Detect if icingaweb2 is installed on the system.
func Detect() bool {
	for _, path := range relevantPaths {
		_, err := os.Stat(path)
		if err == nil {
			return true
		}
	}

	return false
}

// Collect data for icingaweb2.
func Collect(c *collection.Collection) {
	if !Detect() {
		c.Log.Info("Could not find icingaweb2")
		return
	}

	c.Log.Info("Collecting Icinga Web 2 information")

	c.AddInstalledPackagesRaw(ModuleName+"/packages.txt", "*icingaweb2*", "*icingacli*")

	if _, ok := collection.IsGitRepository("/usr/share/icingaweb2"); ok {
		c.AddGitRepoInfo(ModuleName+"/git.yml", "/usr/share/icingaweb2")
	}

	CollectModuleInfo(c)

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for name, cmd := range commands {
		c.AddCommandOutput(ModuleName+"/"+name, cmd[0], cmd[1:]...)
	}

	// Detect PHP related packages and services
	c.AddInstalledPackagesRaw(ModuleName+"/packages-php.txt", "*php*")

	if services, err := collection.FindServices("*php*-fpm"); err == nil && len(services) > 0 {
		for _, name := range services {
			c.AddServiceStatusRaw(ModuleName+"/service-"+name+".txt", name)
		}
	}

	// Detect webserver packages
	c.AddInstalledPackagesRaw(ModuleName+"/packages-webserver.txt", "*apache*", "*httpd*")
}

func CollectModuleInfo(c *collection.Collection) {
	if !collection.DetectGitInstalled() {
		c.Log.Warnf("we need git to inspect modules closer")
	}

	modulesFiles, err := ioutil.ReadDir(ModulesPath)
	if err != nil {
		c.Log.Warnf("Could not list modules in %s - %s", ModulesPath, err)
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
