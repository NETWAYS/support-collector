package base

import (
	"bytes"
	"fmt"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"gopkg.in/yaml.v3"
)

func Collect(c *collection.Collection) {
	c.Log.Info("Collecting base system information")

	CollectKernelInfo(c)
}

func CollectKernelInfo(c *collection.Collection) {
	buf := bytes.Buffer{}

	err := yaml.NewEncoder(&buf).Encode(GetKernelInfo())
	if err != nil {
		// TODO: logging
		fmt.Println(err)
	}

	err = c.AddFileFromReader("kernel.yml", &buf)
	if err != nil {
		// TODO: logging
		fmt.Println(err)
	}
}
