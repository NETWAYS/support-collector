package main

import (
	"fmt"
	"github.com/NETWAYS/support-collector/modules/base"
	"github.com/NETWAYS/support-collector/modules/icinga2"
	"github.com/NETWAYS/support-collector/modules/icingadirector"
	"github.com/NETWAYS/support-collector/modules/icingaweb2"
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
}

var (
	outputFile                      string
	enabledModules, disabledModules []string
	debug, printVersion             bool
)

func handleArguments() {
	// Build default list of enabled modules
	enabledModules = make([]string, 0, len(modules))
	for k := range modules {
		enabledModules = append(enabledModules, k)
	}

	flag.StringVarP(&outputFile, "output", "o", "support-collector.zip", "Output file for the ZIP content")
	flag.StringSliceVar(&enabledModules, "enable", enabledModules, "List of enabled module")
	flag.StringSliceVar(&disabledModules, "disable", []string{}, "List of disabled module")
	flag.BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	flag.BoolVarP(&printVersion, "version", "V", false, "Print version and exit")
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

func main() {
	handleArguments()

	// Prepare output
	file, err := os.Create(outputFile)
	if err != nil {
		logrus.Fatal(err)
	}

	c := collection.New(file)

	if debug {
		c.Log.SetLevel(logrus.DebugLevel)
	}

	// Add console log output via logrus.Hook
	c.Log.AddHook(&util.ExtraLogHook{
		Formatter: &logrus.TextFormatter{ForceColors: true},
		Writer:    colorable.NewColorableStdout(),
	})

	// set locale to C, to avoid translations in command output
	_ = os.Setenv("LANG", "C")

	c.Log.Infof("Starting %s", Product) // TODO: add version

	if !isPrivilegedUser() {
		c.Log.Warn("This tool should be run as a privileged user (root) to collect all necessary information")
	}

	startTime := time.Now()

	// Call all enabled modules
	for name, call := range modules {
		switch {
		case stringInSlice(name, disabledModules):
			c.Log.Infof("Module %s is disabled", name)
		case !stringInSlice(name, enabledModules):
			c.Log.Infof("Module %s is not enabled", name)
		default:
			c.Log.Debugf("Calling module %s", name)
			call(c)
		}
	}

	c.Log.Infof("Collection complete, took us %.3f seconds", time.Since(startTime).Seconds())

	err = c.AddLogToOutput()
	if err != nil {
		logrus.Error(err)
	}

	err = c.Close()
	if err != nil {
		logrus.Error(err)
	}

	err = file.Close()
	if err != nil {
		logrus.Error(err)
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
