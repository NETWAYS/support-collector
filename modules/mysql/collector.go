package mysql

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
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

	c.AddInstalledPackagesRaw(ModuleName+"/packages.txt", "*mysql*", "*mariadb*")
	c.AddServiceStatusRaw(ModuleName+"/service.txt", service)
	c.AddFilesAtLeastOne(ModuleName, possibleConfigPaths...)
}
