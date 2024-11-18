package main

import (
	"encoding/json"
	"fmt"
	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/config"
	"github.com/NETWAYS/support-collector/internal/metrics"
	"github.com/NETWAYS/support-collector/modules/ansible"
	"github.com/NETWAYS/support-collector/modules/base"
	"github.com/NETWAYS/support-collector/modules/corosync"
	"github.com/NETWAYS/support-collector/modules/elastic"
	"github.com/NETWAYS/support-collector/modules/foreman"
	"github.com/NETWAYS/support-collector/modules/grafana"
	"github.com/NETWAYS/support-collector/modules/graphite"
	"github.com/NETWAYS/support-collector/modules/graylog"
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
	"github.com/NETWAYS/support-collector/modules/redis"
	"github.com/NETWAYS/support-collector/modules/webservers"
	flag "github.com/spf13/pflag"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/NETWAYS/support-collector/internal/obfuscate"
	"github.com/NETWAYS/support-collector/internal/util"
	"github.com/NETWAYS/support-collector/modules/icinga2"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

const Product = "NETWAYS support-collector"

const Readme = `
The support-collector allows our customers to collect relevant information from
their servers. A resulting ZIP file can then be provided to our support team
for further inspection.

Find more information and releases at:
		https://github.com/NETWAYS/support-collector

If you are a customer, contact us at:
		support@netways.de  /  https://netways.de/en/contact/

WARNING: DO NOT transfer the generated file over insecure connections or by
email, it contains potential sensitive information!
`

var (
	disableWizard                             bool
	answerFile                                string
	verbose, printVersion, detailedCollection bool
	startTime                                 = time.Now()
)

var modules = map[string]func(*collection.Collection){
	"ansible":         ansible.Collect,
	"base":            base.Collect,
	"corosync":        corosync.Collect,
	"elastic":         elastic.Collect,
	"foreman":         foreman.Collect,
	"grafana":         grafana.Collect,
	"graphite":        graphite.Collect,
	"graylog":         graylog.Collect,
	"icinga-director": icingadirector.Collect,
	"icinga2":         icinga2.Collect,
	"icingadb":        icingadb.Collect,
	"icingaweb2":      icingaweb2.Collect,
	"influxdb":        influxdb.Collect,
	"keepalived":      keepalived.Collect,
	"mongodb":         mongodb.Collect,
	"mysql":           mysql.Collect,
	"postgresql":      postgresql.Collect,
	"prometheus":      prometheus.Collect,
	"puppet":          puppet.Collect,
	"redis":           redis.Collect,
	"webservers":      webservers.Collect,
}

func init() {
	// Set locale to C, to avoid translations in command output
	_ = os.Setenv("LANG", "C")
}

func main() {
	// Create new config object with defaults
	conf := config.GetControlDefaultObject()

	// Add and parse flags
	if err := parseFlags(); err != nil {
		logrus.Fatal(err)
	}

	// Read input from answer-file if provided
	// Needs to done after parsing flags to have the value for answerFile
	if answerFile != "" {
		if err := config.ReadAnswerFile(answerFile, &conf); err != nil {
			logrus.Fatal(err)
		}

		conf.General.AnswerFile = answerFile
	}

	// Start interactive config wizard if not disabled via flag and no answer-file is provided
	if !disableWizard && answerFile == "" {
		startConfigWizard(&conf)
	}

	// If "all" provided for enabled modules, enable all
	if slices.Contains(conf.General.EnabledModules, "all") {
		conf.General.EnabledModules = config.ModulesOrder
	}

	// Validate conf. If errors found, print them and exit
	if validationErrors := config.ValidateConfig(conf); len(validationErrors) > 0 {
		for _, e := range validationErrors {
			logrus.Error(e)
		}

		os.Exit(1)
	}

	// Initialize new collection with default values
	c, closeCollection := NewCollection(conf)

	// Close collection
	defer closeCollection()

	// Initialize new metrics and defer function to save it to json
	c.Metric = metrics.New(getVersion())
	defer func() {
		// Save metrics to file
		body, err := json.Marshal(c.Metric)
		if err != nil {
			c.Log.Warn("cant unmarshal metrics: %w", err)
		}

		c.AddFileJSONRaw("metrics.json", body)
	}()

	c.Metric.Controls = c.Config

	// Choose whether detailed collection will be enabled or not
	if !conf.General.DetailedCollection {
		c.Detailed = false
		c.Config.General.DetailedCollection = false
		c.Log.Warn("Detailed collection is disabled")
	} else {
		c.Detailed = true
		c.Log.Info("Detailed collection is enabled")
	}

	if !util.IsPrivilegedUser() {
		c.Log.Warn("This tool should be run as a privileged user (root) to collect all necessary information")
	}

	// Set command Timeout from argument
	c.ExecTimeout = c.Config.General.CommandTimeout

	// Parse modules
	collectModules(c)

	// Save overall timing
	c.Metric.Timings["total"] = time.Since(startTime)

	c.Log.Infof("Collection complete, took us %.3f seconds", c.Metric.Timings["total"].Seconds())

	// Collect obfuscation info
	var (
		count         uint
		affectedFiles []string
	)

	for _, o := range c.Obfuscators {
		count += o.Replaced

		affectedFiles = append(affectedFiles, o.ObfuscatedFiles...)
	}

	if len(affectedFiles) > 0 {
		c.Log.Infof("Obfuscation replaced %d token in %d files (%d definitions)", count, len(util.DistinctStringSlice(affectedFiles)), len(c.Obfuscators))
	}

	// get absolute path of outputFile
	path, err := filepath.Abs(c.Config.General.OutputFile)
	if err != nil {
		c.Log.Debug(err)
	}

	c.Log.Infof("Generated ZIP file located at %s", path)
}

// NewCollection starts a new collection based on given controls.Config
//
// Collection and cleanup function to defer are returned
func NewCollection(control config.Config) (*collection.Collection, func()) {
	file, err := os.Create(control.General.OutputFile)
	if err != nil {
		logrus.Fatalf("cant create or open file for collection for given value '%s' - %s", control.General.OutputFile, err)
	}

	c := collection.New(file)
	c.Config = control

	consoleLevel := logrus.InfoLevel
	if verbose {
		consoleLevel = logrus.DebugLevel
	}

	c.Log.SetLevel(consoleLevel)

	// Add console log output via logrus.Hook
	c.Log.AddHook(&util.ExtraLogHook{
		Formatter: &logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "15:04:05"},
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

// parseFlags adds the default control arguments and parses them
func parseFlags() (err error) {
	var generateAnswerFile bool
	// General arguments without interactive prompt
	flag.BoolVar(&disableWizard, "disable-wizard", false, "Disable interactive wizard for input via stdin")
	flag.BoolVarP(&printVersion, "version", "v", false, "Print version and exit")
	flag.BoolVarP(&verbose, "verbose", "V", false, "Enable verbose logging")
	flag.BoolVar(&generateAnswerFile, "generate-answer-file", false, "Generate an example answer-file with default values")
	flag.StringVarP(&answerFile, "answer-file", "f", "", "Provide an answer-file to control the collection")

	// Output a proper help message with details
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n%s\n\n", Product, strings.Trim(Readme, "\n"))

		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()
	}

	// Parse flags from command-line
	flag.Parse()

	// Print version and exit
	if printVersion {
		fmt.Println(Product, "version", getBuildInfo()) //nolint:forbidigo
		os.Exit(0)
	}

	if generateAnswerFile {
		err = config.GenerateDefaultAnswerFile()
		if err != nil {
			return err
		}

		os.Exit(0)
	}

	return nil
}

// startConfigWizard will start the stdin config wizard to provide config
func startConfigWizard(conf *config.Config) {
	wizard := config.NewWizard()

	// Define arguments for interactive input via stdin
	wizard.AddStringVar(&conf.General.OutputFile, "output", util.BuildFileName(), "Filename for resulting zip", true, nil)
	wizard.AddSliceVarFromString(&conf.General.EnabledModules, "enable", []string{"all"}, "Which modules should be enabled? (Comma separated list of modules)", false, nil)
	wizard.AddBoolVar(&detailedCollection, "detailed", true, "Enable detailed collection including logs and more (recommended)?", nil)
	wizard.AddStringSliceVar(&conf.General.ExtraObfuscators, "obfuscators", false, "Do you want to define some custom obfuscators (passwords, secrets etc.)", "Add custom obfuscator", nil)

	// Collect Icinga 2 API endpoints if module 'icinga2' is enabled
	// Because we only add this when module 'icinga2' or 'all' is enabled, this needs to be after saving the enabled modules
	wizard.AddIcingaEndpoints(&conf.Icinga2.Endpoints, "icinga-endpoints", "\nModule 'icinga2'is  enabled.\nDo you want to collect data from Icinga 2 API endpoints?", func() bool {
		if ok := util.StringInSlice("all", conf.General.EnabledModules) || slices.Contains(conf.General.EnabledModules, "icinga2"); ok {
			return true
		}

		return false
	})

	wizard.Parse(strings.Join(config.ModulesOrder, ","))

	fmt.Printf("\nArgument wizard finished. Starting...\n\n")
}

func collectModules(c *collection.Collection) {
	// Check if module is enabled / disabled and call it
	for _, name := range config.ModulesOrder {
		switch {
		case util.StringInSlice(name, c.Config.General.DisabledModules):
			c.Log.Debugf("Module %s is disabled", name)
		case !util.StringInSlice(name, c.Config.General.EnabledModules):
			c.Log.Debugf("Module %s is not enabled", name)
		default:
			// Save current time
			moduleStart := time.Now()

			c.Log.Debugf("Start collecting data for module %s", name)

			// Register custom obfuscators
			for _, o := range c.Config.General.ExtraObfuscators {
				c.Log.Debugf("Adding custom obfuscator for '%s' to module %s", o, name)
				c.RegisterObfuscator(obfuscate.NewAny(o))
			}

			// Call collection function for module
			// TODO return errors?
			modules[name](c)

			// Save runtime of module
			c.Metric.Timings[name] = time.Since(moduleStart)

			c.Log.Debugf("Finished with module %s in %.3f seconds", name, c.Metric.Timings[name].Seconds())
		}
	}
}
