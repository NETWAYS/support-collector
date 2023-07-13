package util

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

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

// GetHostnameWithoutDomain returns hostname without domain
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
