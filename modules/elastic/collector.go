package elastic

import (
	"github.com/NETWAYS/support-collector/internal/util"
	"path/filepath"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
)

const ModuleName = "elastic"

var relevantPaths = []string{
	"/etc/elasticsearch/",
	"/etc/logstash/",
	"/etc/kibana/",
}

var files = []string{
	"/etc/elasticsearch/elasticsearch.yaml",
	"/etc/elasticsearch/elasticsearch.yml",
	"/etc/logstash/logstash.yaml",
	"/etc/logstash/logstash.yml",
	"/etc/kibana/kibana.yaml",
	"/etc/kibana/kibana.yml",
}

var possibleDaemons = []string{
	"/usr/lib/systemd/system/elasticsearch.service",
	"/usr/lib/systemd/system/logstash.service",
	"/usr/lib/systemd/system/kibana.service",
}

var services = []string{
	"elasticsearch",
	"kibana",
	"logstash",
}

var obfuscators = []*obfuscate.Obfuscator{
	obfuscate.NewFile(`(?i)(?:password|keypassphrase).*\s*:\s*(.*)`, `yml`),
	obfuscate.NewFile(`(?i)(?:password|keypassphrase).*\s*:\s*(.*)`, `yaml`),
}

func Collect(c *collection.Collection) {
	if !util.ModuleExists(relevantPaths) {
		c.Log.Info("Could not find Elastic Stack")
		return
	}

	c.Log.Info("Collecting Elastic Stack information")

	c.RegisterObfuscators(obfuscators...)

	c.AddInstalledPackagesRaw(filepath.Join(ModuleName, "packages.txt"), "*elastic*", "*logstash*", "*kibana*")

	for _, service := range services {
		c.AddServiceStatusRaw(filepath.Join(ModuleName, "service-"+service+".txt"), service)
	}

	for _, file := range possibleDaemons {
		c.AddFilesIfFound(ModuleName, file)
	}

	for _, file := range files {
		c.AddFilesIfFound(ModuleName, file)
	}
}
