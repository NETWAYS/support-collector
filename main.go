package main

import (
	"fmt"
	"github.com/NETWAYS/support-collector/modules/ansible"
	"github.com/NETWAYS/support-collector/modules/base"
	"github.com/NETWAYS/support-collector/modules/grafana"
	"github.com/NETWAYS/support-collector/modules/graphite"
	"github.com/NETWAYS/support-collector/modules/icinga2"
	"github.com/NETWAYS/support-collector/modules/icingadb"
	"github.com/NETWAYS/support-collector/modules/icingadirector"
	"github.com/NETWAYS/support-collector/modules/icingaweb2"
	"github.com/NETWAYS/support-collector/modules/influxdb"
	"github.com/NETWAYS/support-collector/modules/mysql"
	"github.com/NETWAYS/support-collector/modules/postgresql"
	"github.com/NETWAYS/support-collector/modules/puppet"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/util"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"os"
	"path/filepath"
	"strings"
	"time"
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
    support@netways.de  /  https://netways.de/contact

WARNING: DO NOT transfer the generated file over insecure connections or by
email, it contains potential sensitive information!
`

var modules = map[string]func(*collection.Collection){
	"base":            base.Collect,
	"icinga2":         icinga2.Collect,
	"icingaweb2":      icingaweb2.Collect,
	"icinga-director": icingadirector.Collect,
	"mysql":           mysql.Collect,
	"influxdb":        influxdb.Collect,
	"postgresql":      postgresql.Collect,
	"ansible":         ansible.Collect,
	"puppet":          puppet.Collect,
	"grafana":         grafana.Collect,
	"graphite":        graphite.Collect,
	"icingadb":        icingadb.Collect,
}

var (
	moduleOrder = []string{
		"base",
		"icinga2",
		"icingaweb2",
		"icinga-director",
		"icingadb",
		"mysql",
		"influxdb",
		"postgresql",
		"ansible",
		"puppet",
		"grafana",
		"graphite",
	}
)

var (
	commandTimeout                  = 60 * time.Second
	outputFile                      string
	enabledModules, disabledModules []string
	verbose, printVersion           bool
)

func main() {
	handleArguments()

	// set locale to C, to avoid translations in command output
	_ = os.Setenv("LANG", "C")

	c, cleanup := NewCollection(outputFile)
	defer cleanup()

	if !util.IsPrivilegedUser() {
		c.Log.Warn("This tool should be run as a privileged user (root) to collect all necessary information")
	}

	var (
		startTime = time.Now()
		timings   = map[string]time.Duration{}
	)

	// Set command Timeout from argument
	c.ExecTimeout = commandTimeout

	// Call all enabled modules
	for _, name := range moduleOrder {
		switch {
		case util.StringInSlice(name, disabledModules):
			c.Log.Infof("Module %s is disabled", name)
		case !util.StringInSlice(name, enabledModules):
			c.Log.Infof("Module %s is not enabled", name)
		default:
			moduleStart := time.Now()

			c.Log.Debugf("Calling module %s", name)
			modules[name](c)

			timings[name] = time.Since(moduleStart)
			c.Log.Debugf("Finished with module %s in %.3f seconds", name, timings[name].Seconds())
		}
	}

	// Collect obfuscation info
	var files, count uint

	for _, o := range c.Obfuscators {
		files += o.Files

		count += o.Replaced
	}

	if files > 0 {
		c.Log.Infof("Obfuscation replaced %d token in %d files (%d definitions)", count, files, len(c.Obfuscators))
	}

	// Save timings
	timings["total"] = time.Since(startTime)
	c.Log.Infof("Collection complete, took us %.3f seconds", timings["total"].Seconds())

	c.AddFileYAML("timing.yml", timings)

	path, err := filepath.Abs(outputFile)
	if err != nil {
		c.Log.Debug(err)
	}

	c.Log.Infof("Generated ZIP file located at %s", path)
}

func handleArguments() {
	flag.StringVarP(&outputFile, "output", "o", buildFileName(), "Output file for the ZIP content")
	flag.StringSliceVar(&enabledModules, "enable", moduleOrder, "List of enabled module")
	flag.StringSliceVar(&disabledModules, "disable", []string{}, "List of disabled module")
	flag.DurationVar(&commandTimeout, "command-timeout", commandTimeout, "Timeout for command execution in modules")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	flag.BoolVarP(&verbose, "debug", "d", false, "Enable debug logging (use verbose)")
	flag.BoolVarP(&printVersion, "version", "V", false, "Print version and exit")

	_ = flag.CommandLine.MarkHidden("debug")
	flag.CommandLine.SortFlags = false

	// Output a proper help message with details
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n%s\n\n", Product, strings.Trim(Readme, "\n"))

		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()
	}

	flag.Parse()

	if printVersion {
		fmt.Println(Product, "version", buildVersion()) // nolint:forbidigo
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
	return GetHostnameWithoutDomain() + "-" + FilePrefix + "-" + time.Now().Format("20060102-1504") + ".zip"
}

func GetHostnameWithoutDomain() string {
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Error(err)
	}

	result, _, found := strings.Cut(hostname, ".")
	if !found {
		return hostname
	}

	return result
}

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

	versionString := buildVersion()
	c.Log.Infof("Starting %s version %s", Product, versionString)
	c.AddFileDataRaw("version", []byte(versionString+"\n"))

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
