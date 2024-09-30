package util

import (
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// FilePrefix for the outfile file.
const FilePrefix = "support-collector"

// StringInSlice matches if a string is contained in a slice.
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

// IsPrivilegedUser returns true when the current user is root.
func IsPrivilegedUser() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}

	// TODO: only works on *NIX systems
	return u.Uid == "0"
}

// CommandExists returns true if command exists.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// GetHostnameWithoutDomain returns hostname without domain.
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

// BuildFileName returns a filename to store the output of support collector.
func BuildFileName() string {
	return FilePrefix + "_" + GetHostnameWithoutDomain() + "_" + time.Now().Format("20060102-1504") + ".zip"
}
