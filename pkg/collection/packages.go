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
		if _, err := exec.LookPath(name); err == nil {
			FoundPackageManager = name
			return name
		}
	}

	return ""
}

func ListInstalledPackagesRaw(pattern ...string) ([]byte, error) {
	switch DetectPackageManager() {
	case PackageManagerRPM:
		arguments := []string{"-qa", "--queryformat", rpmQueryFormat}
		arguments = append(arguments, pattern...)

		return LoadCommandOutput("rpm", arguments...)
	case PackageManagerDebian:
		arguments := []string{"-W", "-f", dpkgQueryFormat}
		arguments = append(arguments, pattern...)

		return LoadCommandOutput("dpkg-query", arguments...)
	default:
		return []byte{}, ErrNoPackageManager
	}
}
