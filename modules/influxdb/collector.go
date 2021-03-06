package influxdb

import (
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
)

const ModuleName = "influxdb"

var relevantPaths = []string{
	"/etc/influxdb",
	"/var/lib/influxdb",
}

var files = []string{
	"/etc/influxdb",
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

	c.AddInstalledPackagesRaw(ModuleName+"/packages.txt", "*influx*")
	c.AddServiceStatusRaw(ModuleName+"/service.txt", "influxdb")

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}
}
