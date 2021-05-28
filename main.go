package main

import (
	"fmt"
	"github.com/NETWAYS/support-collector/modules/base"
	"github.com/NETWAYS/support-collector/modules/icinga2"
	"github.com/NETWAYS/support-collector/modules/icingadirector"
	"github.com/NETWAYS/support-collector/modules/icingaweb2"
	"github.com/NETWAYS/support-collector/modules/mysql"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/util"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"os"
	"os/user"
	"time"
)

const Product = "NETWAYS support collector"

var modules = map[string]func(*collection.Collection){
	"base":            base.Collect,
	"icinga2":         icinga2.Collect,
	"icingaweb2":      icingaweb2.Collect,
	"icinga-director": icingadirector.Collect,
	"mysql":           mysql.Collect,
}

var moduleOrder = []string{
	"base",
	"icinga2",
	"icingaweb2",
	"icinga-director",
	"mysql",
}

var (
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

	if !isPrivilegedUser() {
		c.Log.Warn("This tool should be run as a privileged user (root) to collect all necessary information")
	}

	var (
		startTime = time.Now()
		timings   = map[string]time.Duration{}
	)

	// Call all enabled modules
	for _, name := range moduleOrder {
		switch {
		case stringInSlice(name, disabledModules):
			c.Log.Infof("Module %s is disabled", name)
		case !stringInSlice(name, enabledModules):
			c.Log.Infof("Module %s is not enabled", name)
		default:
			moduleStart := time.Now()

			c.Log.Debugf("Calling module %s", name)
			modules[name](c)

			timings[name] = time.Since(moduleStart)
			c.Log.Debugf("Finished with module %s in %.3f seconds", name, timings[name].Seconds())
		}
	}

	timings["total"] = time.Since(startTime)
	c.Log.Infof("Collection complete, took us %.3f seconds", timings["total"].Seconds())

	c.AddFileYAML("timing.yml", timings)
}

func handleArguments() {
	flag.StringVarP(&outputFile, "output", "o", "support-collector.zip", "Output file for the ZIP content")
	flag.StringSliceVar(&enabledModules, "enable", moduleOrder, "List of enabled module")
	flag.StringSliceVar(&disabledModules, "disable", []string{}, "List of disabled module")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	flag.BoolVarP(&verbose, "debug", "d", false, "Enable debug logging (use verbose)")
	flag.BoolVarP(&printVersion, "version", "V", false, "Print version and exit")

	_ = flag.CommandLine.MarkHidden("debug")
	flag.CommandLine.SortFlags = false

	// TODO: Add usage with some documentation
	/*
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		}
	*/

	flag.Parse()

	if printVersion {
		fmt.Println(Product, " version ", buildVersion()) // nolint:forbidigo
		os.Exit(0)
	}

	// Verify enabled modules
	for _, name := range enabledModules {
		if _, ok := modules[name]; !ok {
			logrus.Fatal("Unknown module to enable: ", name)
		}
	}
}

func isPrivilegedUser() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}

	// TODO: only works on *NIX systems
	return u.Uid == "0"
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
	c.AddFileData("version", []byte(versionString+"\n"))

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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}
