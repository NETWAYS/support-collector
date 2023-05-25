package webservers

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
	"path/filepath"
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

func DetectWebservers() bool {
	for _, path := range relevantPaths {
		_, err := os.Stat(path)
		if err == nil {
			return true
		}
	}

	return false
}

func Collect(c *collection.Collection) {
	if !DetectWebservers() {
		c.Log.Info("Could not find webservers")
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
}
