package collection

import (
	"errors"
	"os/exec"
)

const (
	PackageManagerRPM    = "rpm"
	PackageManagerDebian = "dpkg"
)

const (
	rpmQueryFormat  = `%{NAME} %{VERSION}-%{RELEASE}\n`
	dpkgQueryFormat = `${Package} ${Version} ${Architecture} ${Status}\n`
)

// ErrNoPackageManager is returned when we could not detect one.
var ErrNoPackageManager = errors.New("could not detect a supported package manager")

var FoundPackageManager string

func DetectPackageManager() string {
	if FoundPackageManager != "" {
		return FoundPackageManager
	}

	priority := []string{PackageManagerDebian, PackageManagerRPM}

	for _, name := range priority {
		_, err := exec.LookPath(name)
		if err == nil {
			FoundPackageManager = name
			return name
		}
	}

	return ""
}

func ListInstalledPackagesRaw(pattern ...string) ([]byte, error) {
	switch DetectPackageManager() {
	case PackageManagerRPM:
		arguments := make([]string, 3+len(pattern))
		arguments[0] = "-qa"
		arguments[1] = "--queryformat"
		arguments[2] = rpmQueryFormat

		for i := range pattern {
			arguments[i+3] = pattern[i]
		}

		return LoadCommandOutput("rpm", arguments...)
	case PackageManagerDebian:
		arguments := make([]string, 3+len(pattern))
		arguments[0] = "-W"
		arguments[1] = "-f"
		arguments[2] = dpkgQueryFormat

		for i := range pattern {
			arguments[i+3] = pattern[i]
		}

		return LoadCommandOutput("dpkg-query", arguments...)
	default:
		return []byte{}, ErrNoPackageManager
	}
}
