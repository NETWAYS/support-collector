package influxdb

import (
	"os"
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
)

const ModuleName = "influxdb"

var relevantPaths = []string{
	"/etc/influxdb",
	"/var/lib/influxdb",
}

var files = []string{
	"/etc/influxdb",
}

var detailedFiles = []string{
	"/var/log/influxdb",
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
		c.Log.Info("Could not find InfluxDB")
		return
	}

	c.Log.Info("Collecting InfluxDB information")

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"), "*influx*")
	c.AddServiceStatusRaw(filepath.Join(ModuleName, "service.txt"), "influxdb")

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	if c.Detailed {
		for _, file := range detailedFiles {
			c.AddFilesIfFound(ModuleName, file)
		}
	}
}
