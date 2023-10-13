package hosts

import (
	"encoding/json"
	"github.com/NETWAYS/support-collector/pkg/collection"
	"os"
)

type HostList struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	Hostname  string `json:"hostname"`
	Localhost bool   `json:"localhost"`
	IpAddr    string `json:"ip_addr"`
	Port      string `json:"port"`
	AuthType  string `json:"auth_type"`
	User      string `json:"user"`
	Password  string `json:"password"`
	Desc      string `json:"description"`
}

// ReadFromJson reads hosts from given json file
func ReadFromJson(hostsFile string) ([]Host, error) {
	file, err := os.ReadFile(hostsFile)
	if err != nil {
		return []Host{}, err
	}

	var hosts HostList
	err = json.Unmarshal(file, &hosts)
	if err != nil {
		return []Host{}, err
	}

	return hosts.Hosts, nil
}

// Collect collects data from the host
func (host *Host) Collect(c *collection.Collection) error {
	switch host.Localhost {
	case true:
		err := host.collectLocal(c)
		if err != nil {
			return err
		}
		// TODO collect local data
		break
	case false:
		err := host.collectRemote(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (host *Host) collectLocal(c *collection.Collection) error {
	// TODO collect local data
	return nil
}

func (host *Host) collectRemote(c *collection.Collection) error {
	err := host.Prepare()
	if err != nil {
		return err
	}
	// TODO collect remote data
	// https://pkg.go.dev/github.com/melbahja/goph#section-readme

	return nil
}

// IsLocalhost returns true if host is localhost
func (host *Host) IsLocalhost() bool {
	if !host.Localhost {
		return false
	}
	return true
}
