package main

import (
	"fmt"
	"github.com/NETWAYS/support-collector/pkg/base"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/util"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"os"
)

const Product = "NETWAYS support collector"

var modules = map[string]func(*collection.Collection){
	"base": base.Collect,
}

var (
	enabledModules, disabledModules []string
	debug, version                  bool
)

func init() {
	enabledModules = make([]string, 0, len(modules))
	for k := range modules {
		enabledModules = append(enabledModules, k)
	}
}

func handleArguments(c *collection.Collection) {
	flag.StringSliceVar(&enabledModules, "enable", enabledModules, "List of enabled module")
	flag.StringSliceVar(&disabledModules, "disable", []string{}, "List of enabled module")
	flag.BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	flag.BoolVarP(&version, "version", "V", false, "Print version and exit")
	flag.CommandLine.SortFlags = false

	// TODO: Add usage with some documentation
	/*
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
		}
	*/

	flag.Parse()

	if version {
		// TODO: print version
		fmt.Println(Product)
		os.Exit(0)
	}

	if debug {
		c.Log.SetLevel(logrus.DebugLevel)
	}

	// Verify enabled modules
	for _, name := range enabledModules {
		if _, ok := modules[name]; !ok {
			fmt.Println("Unknown module to enable:", name)
			os.Exit(1)
		}
	}
}

func main() {
	c := collection.New()

	handleArguments(c)

	// Add console log output via logrus.Hook
	c.Log.AddHook(&util.ExtraLogHook{
		Formatter: &logrus.TextFormatter{ForceColors: true},
		Writer:    colorable.NewColorableStdout(),
	})

	c.Log.Infof("Starting %s", Product) // TODO: add version

	// Call all enabled modules
	for name, call := range modules {
		if stringInSlice(name, disabledModules) {
			c.Log.Infof("Module %s is disabled", name)
		} else if !stringInSlice(name, enabledModules) {
			c.Log.Infof("Module %s is not enabled", name)
		} else {
			c.Log.Infof("Calling module %s", name)
			call(c)
		}
	}

	// Write out the ZIP file
	file, err := os.Create("support-collector.zip")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = c.WriteZIP(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file.Close()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
