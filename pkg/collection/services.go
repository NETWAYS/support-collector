package collection

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	ServiceManagerSystemD = "systemd"
	ServiceManagerSysV    = "sysv"
)

const sysVInitD = "/etc/init.d"

// ErrNoServiceManager is returned when we could not detect one.
var ErrNoServiceManager = errors.New("could not detect a supported service manager")

// FoundServiceManager remembers the current service manager found.
var FoundServiceManager string

var systemdUnitPaths = []string{
	"/etc/systemd/system",
	"/usr/lib/systemd/system",
	"/lib/systemd/system",
}

func DetectServiceManager() string {
	if FoundServiceManager != "" {
		return FoundServiceManager
	}

	priority := map[string]string{
		"systemctl": ServiceManagerSystemD,
		"service":   ServiceManagerSysV,
	}

	for command, manager := range priority {
		if _, err := exec.LookPath(command); err == nil {
			FoundServiceManager = manager
			return manager
		}
	}

	return ""
}

func GetServiceStatusRaw(name string) ([]byte, error) {
	switch DetectServiceManager() {
	case ServiceManagerSystemD:
		return LoadCommandOutput("systemctl", "status", "--full", name+".service")
	case ServiceManagerSysV:
		return LoadCommandOutput("service", name, "status")
	default:
		return []byte{}, ErrNoServiceManager
	}
}

func FindServices(pattern string) (map[string]string, error) {
	switch DetectServiceManager() {
	case ServiceManagerSystemD:
		return FindServicesSystemd(pattern)
	case ServiceManagerSysV:
		return FindServicesSysV(pattern)
	default:
		return nil, ErrNoServiceManager
	}
}

func FindServicesSystemd(pattern string) (map[string]string, error) {
	units := map[string]string{}

	for _, path := range systemdUnitPaths {
		files, err := filepath.Glob(filepath.Join(path, pattern))
		if err != nil {
			return nil, fmt.Errorf("could not glob with pattern '%s': %w", pattern, err)
		}

		for _, file := range files {
			// Skip the file if it is a symlink
			if stat, err := os.Lstat(file); err != nil || (stat.Mode()&os.ModeSymlink) == os.ModeSymlink {
				continue
			}

			name := filepath.Base(file)

			// Only safe the first found unit
			if _, ok := units[name]; !ok {
				units[name] = file
			}
		}
	}

	return units, nil
}

func FindServicesSysV(pattern string) (map[string]string, error) {
	units := map[string]string{}

	files, err := filepath.Glob(filepath.Join(sysVInitD, pattern))
	if err != nil {
		return nil, fmt.Errorf("could not glob with pattern '%s': %w", pattern, err)
	}

	for _, file := range files {
		units[filepath.Base(file)] = file
	}

	return units, nil
}
