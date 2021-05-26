package main

import (
	"fmt"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"github.com/NETWAYS/support-collector/pkg/util"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	c := collection.New()

	// Add console log output via logrus.Hook
	c.Log.AddHook(&util.ExtraLogHook{
		Formatter: &logrus.TextFormatter{ForceColors: true},
		Writer:    colorable.NewColorableStdout(),
	})

	c.Log.Info("Starting NETWAYS support collector")

	err := c.AddFiles("test", "pkg/collection/testdata")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = c.AddCommandOutput("test/output.txt", "sh", "-c", "echo testoutput")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
