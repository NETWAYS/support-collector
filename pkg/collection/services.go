package collection

import (
	"fmt"
	"os/exec"
)

const (
	ServiceManagerSystemD = "systemd"
	ServiceManagerSysV    = "sysv"
)

var FoundServiceManager string

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
		return []byte{}, fmt.Errorf("could not detect a supported service manager")
	}
}
