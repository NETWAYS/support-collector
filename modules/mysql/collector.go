package mysql

import (
	"os"
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
)

const (
	ModuleName = "mysql"
)

// Possible services for MySQL.
var possibleServices = []string{
	"mysql",
	"mysqld",
	"mariadb",
}

// Possible config paths to collect, the first paths can not be a glob.
var possibleConfigPaths = []string{
	"/etc/my.cnf*",
	"/etc/mysql*",
	"/etc/mariadb*",
}

var commands = map[string][]string{
	"mysql-version.txt": {"mysql", "-V"},
}

var optionalFiles = []string{
	"/etc/logrotate.d/mariadb",
	"/etc/logrotate.d/mysql",
}

var detailedFiles = []string{
	"/var/log/mysql/error.log",
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:password)=(.*)`, `cnf`),
}

// Detect if a MySQL or MariaDB daemon appears to be running.
func Detect() string {
	for _, name := range possibleServices {
		_, err := collection.GetServiceStatusRaw(name)
		if err == nil {
			return name
		}
	}

	return ""
}

// Collect data for MySQL or MariaDB.
func Collect(c *collection.Collection) {
	service := Detect()
	if service == "" {
		c.Log.Info("Could not a running MySQL or MariaDB service")
		return
	}

	c.Log.Info("Collecting MySQL/MariaDB information")

	c.RegisterObfuscators(obfuscators...)

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"), "*mysql*", "*mariadb*")
	c.AddServiceStatusRaw(filepath.Join(ModuleName, "service.txt"), service)
	c.AddFilesIfFound(ModuleName, possibleConfigPaths...)

	for name, cmd := range commands {
		c.AddCommandOutput(filepath.Join(ModuleName, name), cmd[0], cmd[1:]...)
	}

	for _, file := range optionalFiles {
		if _, err := os.Stat(file); err != nil {
			continue
		}

		c.AddFiles(ModuleName, file)
	}

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFilesIfFound(ModuleName, file)
		}
	}
}
