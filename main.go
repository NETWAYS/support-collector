package main

import (
	"fmt"
	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/collection/modules"
	"github.com/NETWAYS/support-collector/internal/connector"
	util2 "github.com/NETWAYS/support-collector/internal/util"
	"github.com/mattn/go-colorable"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"os"
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

var (
	verbose, printVersion           bool
	enabledModules, disabledModules []string
	extraObfuscators                []string
	outputFile                      string
	commandTimeout                  = 60 * time.Second
	noDetailedCollection            bool
	hostsFile                       string
)

func main() {
	handleArguments()

	// set locale to C, to avoid translations in command output
	_ = os.Setenv("LANG", "C")

	// initialize new collection
	c, cleanup := NewCollection(outputFile)
	defer cleanup()

	// CollectHosts hosts from given json
	if hostsFile != "" {
		h, err := connector.CollectHosts(hostsFile)
		if err != nil {
			c.Log.Fatalf("cant read hosts from json, %s", err)
		}

		for _, host := range h {
			c.Log.Infof("Start collection for %s", host.Hostname)

			err = host.Collect(c)
			if err != nil {
				c.Log.Warn(err)
			}
			c.Log.Infof("Finished collection for %s", host.Hostname)
		}
	}
}

func handleArguments() {
	flag.StringVarP(&outputFile, "output", "o", buildFileName(), "Output file for the ZIP content")
	flag.StringVar(&hostsFile, "hosts", "", "Path to hosts file")
	flag.StringSliceVar(&enabledModules, "enable", modules.Order, "List of enabled module")
	flag.StringSliceVar(&disabledModules, "disable", []string{}, "List of disabled module")
	flag.BoolVar(&noDetailedCollection, "nodetails", false, "Disable detailed collection including logs and more")
	flag.StringArrayVar(&extraObfuscators, "hide", []string{}, "List of additional strings to obfuscate. Can be used multiple times and supports regex.") //nolint:lll
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
		fmt.Println(Product, "version", buildVersion()) //nolint:forbidigo
		os.Exit(0)
	}

	// Verify enabled modules
	for _, name := range enabledModules {
		if _, ok := modules.List[name]; !ok {
			logrus.Fatal("Unknown module to enable: ", name)
		}
	}
}

// buildFileName returns a filename to store the output of support collector.
func buildFileName() string {
	return util2.GetHostnameWithoutDomain() + "-" + FilePrefix + "-" + time.Now().Format("20060102-1504") + ".zip"
}

// NewCollection initializes a new collection
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
	c.Log.AddHook(&util2.ExtraLogHook{
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
