package main

import (
	"encoding/json"
	"fmt"
	"github.com/NETWAYS/support-collector/internal/metrics"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/obfuscate"
	"github.com/NETWAYS/support-collector/internal/util"
	"github.com/NETWAYS/support-collector/modules/ansible"
	"github.com/NETWAYS/support-collector/modules/base"
	"github.com/NETWAYS/support-collector/modules/corosync"
	"github.com/NETWAYS/support-collector/modules/elastic"
	"github.com/NETWAYS/support-collector/modules/foreman"
	"github.com/NETWAYS/support-collector/modules/grafana"
	"github.com/NETWAYS/support-collector/modules/graphite"
	"github.com/NETWAYS/support-collector/modules/graylog"
	"github.com/NETWAYS/support-collector/modules/icinga2"
	"github.com/NETWAYS/support-collector/modules/icingadb"
	"github.com/NETWAYS/support-collector/modules/icingadirector"
	"github.com/NETWAYS/support-collector/modules/icingaweb2"
	"github.com/NETWAYS/support-collector/modules/influxdb"
	"github.com/NETWAYS/support-collector/modules/keepalived"
	"github.com/NETWAYS/support-collector/modules/mongodb"
	"github.com/NETWAYS/support-collector/modules/mysql"
	"github.com/NETWAYS/support-collector/modules/postgresql"
	"github.com/NETWAYS/support-collector/modules/prometheus"
	"github.com/NETWAYS/support-collector/modules/puppet"
	"github.com/NETWAYS/support-collector/modules/webservers"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const Product = "NETWAYS support collector"

// FilePrefix for the outfile file.
const FilePrefix = "netways-support"

const Readme = `
The support collector allows our customers to collect relevant information from
their servers. A resulting ZIP file can then be provided to our support team
for further inspection.

Find more information and releases at:
		https://github.com/NETWAYS/support-collector

If you are a customer, contact us at:
		support@netways.de  /  https://netways.de/en/contact/

WARNING: DO NOT transfer the generated file over insecure connections or by
email, it contains potential sensitive information!
`

var modules = map[string]func(*collection.Collection){
	"base":            base.Collect,
	"webservers":      webservers.Collect,
	"icinga2":         icinga2.Collect,
	"icingaweb2":      icingaweb2.Collect,
	"icinga-director": icingadirector.Collect,
	"elastic":         elastic.Collect,
	"corosync":        corosync.Collect,
	"keepalived":      keepalived.Collect,
	"mongodb":         mongodb.Collect,
	"mysql":           mysql.Collect,
	"influxdb":        influxdb.Collect,
	"postgresql":      postgresql.Collect,
	"prometheus":      prometheus.Collect,
	"ansible":         ansible.Collect,
	"puppet":          puppet.Collect,
	"grafana":         grafana.Collect,
	"graphite":        graphite.Collect,
	"graylog":         graylog.Collect,
	"icingadb":        icingadb.Collect,
	"foreman":         foreman.Collect,
}

var (
	moduleOrder = []string{
		"base",
		"webservers",
		"icinga2",
		"icingaweb2",
		"icinga-director",
		"icingadb",
		"elastic",
		"corosync",
		"keepalived",
		"mongodb",
		"mysql",
		"influxdb",
		"postgresql",
		"prometheus",
		"ansible",
		"puppet",
		"grafana",
		"graphite",
		"graylog",
		"foreman",
	}
)

var (
	verbose, printVersion           bool
	enabledModules, disabledModules []string
	extraObfuscators                []string
	outputFile                      string
	commandTimeout                  = 60 * time.Second
	noDetailedCollection            bool
	startTime                       = time.Now()
	metric                          *metrics.Metrics
)

func main() {
	handleArguments()

	// Set locale to C, to avoid translations in command output
	_ = os.Setenv("LANG", "C")

	c, closeCollection := NewCollection(outputFile)
	// Close collection
	defer closeCollection()

	// Initialize new metrics and defer function to save it to json
	metric = metrics.New(getVersion())
	defer func() {
		// Save metrics to file
		body, err := json.Marshal(metric)
		if err != nil {
			c.Log.Warn("cant unmarshal metrics: %w", err)
		}

		c.AddFileJSON("metrics.json", body)
	}()

	if noDetailedCollection {
		c.Detailed = false
		c.Log.Warn("Detailed collection is disabled")
	} else {
		c.Log.Info("Detailed collection is enabled")
	}

	if !util.IsPrivilegedUser() {
		c.Log.Warn("This tool should be run as a privileged user (root) to collect all necessary information")
	}

	// Set command Timeout from argument
	c.ExecTimeout = commandTimeout

	// Collect modules
	collectModules(c)

	// Save overall timing
	metric.Timings["total"] = time.Since(startTime)

	c.Log.Infof("Collection complete, took us %.3f seconds", metric.Timings["total"].Seconds())

	// Collect obfuscation info
	var files, count uint

	for _, o := range c.Obfuscators {
		files += o.Files

		count += o.Replaced
	}

	if files > 0 {
		c.Log.Infof("Obfuscation replaced %d token in %d files (%d definitions)", count, files, len(c.Obfuscators))
	}

	// get absolute path of outputFile
	path, err := filepath.Abs(outputFile)
	if err != nil {
		c.Log.Debug(err)
	}

	c.Log.Infof("Generated ZIP file located at %s", path)
}

func handleArguments() {
	// arguments for collection handling
	flag.StringSliceVar(&enabledModules, "enable", moduleOrder, "List of enabled module")
	flag.StringSliceVar(&disabledModules, "disable", []string{}, "List of disabled module")
	flag.StringVarP(&outputFile, "output", "o", buildFileName(), "Output file for the ZIP content")
	flag.BoolVar(&noDetailedCollection, "nodetails", false, "Disable detailed collection including logs and more")
	flag.StringArrayVar(&extraObfuscators, "hide", []string{}, "List of additional strings to obfuscate. Can be used multiple times and supports regex.") //nolint:lll
	flag.DurationVar(&commandTimeout, "command-timeout", commandTimeout, "Timeout for command execution in modules")

	// api credentials for icinga 2 modules
	flag.StringVar(&icinga2.APICred.Username, "icinga2-api-user", "", "Username of global Icinga 2 API user to collect data about Icinga 2 Infrastructure")                                                                          //nolint:lll
	flag.StringVar(&icinga2.APICred.Password, "icinga2-api-pass", "", "Password for global Icinga 2 API user to collect data about Icinga 2 Infrastructure")                                                                         //nolint:lll
	flag.StringSliceVar(&icinga2.APIEndpoints, "icinga2-api-endpoints", []string{}, "List of Icinga 2 API Endpoints (including port) to collect data from. FQDN or IP address must be reachable. (Example: i2-master01.local:5665)") //nolint:lll

	// basic arguments
	flag.BoolVarP(&printVersion, "version", "V", false, "Print version and exit")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")

	flag.CommandLine.SortFlags = false

	// Output a proper help message with details
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n%s\n\n", Product, strings.Trim(Readme, "\n"))

		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()
	}

	flag.Parse()

	if printVersion {
		fmt.Println(Product, "version", getBuildInfo()) //nolint:forbidigo
		os.Exit(0)
	}

	// Verify enabled modules
	for _, name := range enabledModules {
		if _, ok := modules[name]; !ok {
			logrus.Fatal("Unknown module to enable: ", name)
		}
	}
}

// buildFileName returns a filename to store the output of support collector.
func buildFileName() string {
	return util.GetHostnameWithoutDomain() + "-" + FilePrefix + "-" + time.Now().Format("20060102-1504") + ".zip"
}

// NewCollection starts a new collection. outputFile will be created.
//
// Collection and cleanup function to defer are returned
func NewCollection(outputFile string) (*collection.Collection, func()) {
	file, err := os.Create(outputFile)
	if err != nil {
		logrus.Fatal(err)
	}

	c := collection.New(file)
	c.Log.SetLevel(logrus.DebugLevel)

	consoleLevel := logrus.InfoLevel
	if verbose {
		// logrus.StandardLogger().SetLevel(logrus.DebugLevel)
		consoleLevel = logrus.DebugLevel
	}

	// Add console log output via logrus.Hook
	c.Log.AddHook(&util.ExtraLogHook{
		Formatter: &logrus.TextFormatter{ForceColors: true},
		Writer:    colorable.NewColorableStdout(),
		Level:     consoleLevel,
	})

	versionString := getBuildInfo()
	c.Log.Infof("Starting %s version %s", Product, versionString)

	return c, func() {
		// Close all open outputs in order, but only log errors
		for _, call := range []func() error{
			c.AddLogToOutput,
			c.Close,
			file.Close,
		} {
			err = call()
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}

func collectModules(c *collection.Collection) {
	// Check if module is enabled / disabled and call it
	for _, name := range moduleOrder {
		switch {
		case util.StringInSlice(name, disabledModules):
			c.Log.Debugf("Module %s is disabled", name)
		case !util.StringInSlice(name, enabledModules):
			c.Log.Debugf("Module %s is not enabled", name)
		default:
			// Save current time
			moduleStart := time.Now()

			c.Log.Debugf("Start collecting data for module %s", name)

			// Register custom obfuscators
			for _, o := range extraObfuscators {
				c.Log.Debugf("Adding custom obfuscator for '%s' to module %s", o, name)
				c.RegisterObfuscator(obfuscate.NewAny(o))
			}

			// Call collection function for module
			modules[name](c)

			// Save runtime of module
			metric.Timings[name] = time.Since(moduleStart)

			c.Log.Debugf("Finished with module %s in %.3f seconds", name, metric.Timings[name].Seconds())
		}
	}
}
