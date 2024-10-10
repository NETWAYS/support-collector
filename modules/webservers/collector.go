package webservers

import (
	"github.com/NETWAYS/support-collector/internal/util"
	"os"
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
)

const ModuleName = "webserver"

var relevantPaths = []string{
	"/etc/apache2",
	"/etc/httpd",
	"/etc/nginx",
}

var optionalFiles = []string{
	"/etc/apache2",
	"/etc/logrotate.d/apache2",
	"/etc/httpd",
	"/etc/logrotate.d/httpd",
	"/etc/nginx",
	"/etc/logrotate.d/nginx",
}

var detailedFiles = []string{
	"/var/log/apache2",
	"/var/log/httpd",
	"/var/log/nginx",
}

var services = []string{
	"apache2",
	"nginx",
	"httpd",
}

var possibleDaemons = []string{
	"/lib/systemd/system/apache2.service",
	"/usr/lib/systemd/system/nginx.service",
	"/lib/systemd/system/nginx.service",
	"/usr/lib/systemd/system/httpd.service",
}

func Collect(c *collection.Collection) {
	if !util.ModuleExists(relevantPaths) {
		c.Log.Info("Could not find webservers. Skipping")
		return
	}

	c.Log.Info("Collecting webservers information")

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"),
		"apache2",
		"nginx",
		"httpd",
	)

	for _, file := range optionalFiles {
		if _, err := os.Stat(file); err != nil {
			continue
		}

		c.AddFiles(ModuleName, file)
	}

	for _, service := range services {
		c.AddServiceStatusRaw(filepath.Join(ModuleName, "service-"+service+".txt"), service)
	}

	for _, file := range possibleDaemons {
		if _, err := os.Stat(file); err != nil {
			continue
		}

		c.AddFilesIfFound(ModuleName, file)
	}

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFilesIfFound(ModuleName, file)
		}
	}
}
