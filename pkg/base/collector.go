package base

import (
	"bytes"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"gopkg.in/yaml.v3"
	"os/exec"
)

const ModuleName = "base"

var files = []string{
	"/etc/os-release",
	"/proc/cpuinfo",
	"/proc/meminfo",
	"/proc/loadavg",
}

var commands = [][]string{
	{"lsblk"},
	{"lspci"},
	{"lsusb"},
	{"dmidecode"},
	{"df", "-T"},
	{"top", "-b", "-n1"},
}

func Collect(c *collection.Collection) {
	c.Log.Info("Collecting base system information")

	CollectKernelInfo(c)

	// Check if apparmor is installed and get status
	if _, err := exec.LookPath("apparmor_status"); err == nil {
		c.AddCommandOutput(ModuleName+"apparmor-status.txt", "apparmor_status")
	}

	// Check if we can detect SELinux enforcing
	if _, err := exec.LookPath("getenforce"); err == nil {
		c.AddCommandOutput(ModuleName+"selinux-status.txt", "getenforce")
	}

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	for _, cmd := range commands {
		name := ModuleName + "/" + cmd[0] + ".txt"
		c.AddCommandOutput(name, cmd[0], cmd[1:]...)
	}
}

func CollectKernelInfo(c *collection.Collection) {
	buf := bytes.Buffer{}

	c.Log.Debug("Collecting Kernel and OS infos")

	info, err := GetKernelInfo()
	if err != nil {
		c.Log.Error(err)
	}

	err = yaml.NewEncoder(&buf).Encode(info)
	if err != nil {
		c.Log.Error(err)
		return
	}

	err = c.AddFileFromReader("kernel.yml", &buf)
	if err != nil {
		c.Log.Error(err)
	}
}
