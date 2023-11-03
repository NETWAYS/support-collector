package base

import (
	"fmt"
	"github.com/NETWAYS/support-collector/internal/collection"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
)

const ModuleName = "base"

var files = []string{
	"/etc/os-release",
	"/proc/cpuinfo",
	"/proc/meminfo",
	"/proc/loadavg",
}

var repositoryFiles = []string{
	"/etc/apt/sources.list",
	"/etc/apt/sources.list.d/",
	"/etc/yum.repos.d/",
	"/etc/zypp/repos.d/",
}

var commands = [][]string{
	{"top", "-b", "-n3"},
	/*{"lsblk"},
	{"lspci"},
	{"lsusb"},
	{"dmidecode"},
	{"df", "-T"},
	{"iotop", "-b", "-n3"},
	{"ioping", "/dev/sda", "-c5"},
	*/
}

func CollectLocal(c *collection.Collection) {
	c.Log.Info("Collecting base system information")

	//CollectKernelInfo(c) TODO add again

	for _, cmd := range commands {
		c.AddCommandOutput(filepath.Join(ModuleName, cmd[0]+".txt"), cmd[0], cmd[1:]...)
	}

	/* TODO add again
	// Check if apparmor is installed and get status
	if _, err := exec.LookPath("apparmor_status"); err == nil {
		c.AddCommandOutput(filepath.Join(ModuleName, "apparmor-status.txt"), "apparmor_status")
	}

	// Check if we can detect SELinux enforcing
	for _, cmd := range []string{"sestatus", "getenforce"} {
		if _, err := exec.LookPath(cmd); err == nil {
			c.AddCommandOutput(filepath.Join(ModuleName, "selinux-status.txt"), cmd)
			break
		}
	}

	for _, file := range files {
		c.AddFiles(ModuleName, file)
	}

	// Add repository settings, at least one of the locations should be found
	c.AddFilesIfFound(ModuleName, repositoryFiles...)
	*/
}

func CollectRemote(client *ssh.Client, c *collection.Collection) {
	for _, cmd := range []string{"lsblk", "ls -al", "lsusb"} {

		file, _ := os.Create(cmd + ".out")

		session, err := client.NewSession()
		if err != nil {
			c.Log.Fatal(err)
		}

		session.Stdout = file

		// Start the command
		if err := session.Start(cmd); err != nil {
			c.Log.Fatalf("Failed to start SSH session: %v", err)
		}
		if err := session.Wait(); err != nil {
			c.Log.Fatal(fmt.Errorf("error waiting: %w", err))
		}
		session.Close()
	}

	/* TODO add again
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

	   	err = c.AddFileFromReaderRaw(filepath.Join(ModuleName, "kernel.yml"), &buf)
	   	if err != nil {
	   		c.Log.Error(err)
	   	}
	   }
	*/
}
