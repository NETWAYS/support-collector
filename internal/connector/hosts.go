package connector

import (
	"encoding/json"
	"github.com/NETWAYS/support-collector/internal/collection"
	"github.com/NETWAYS/support-collector/internal/collection/modules/base"
	"os"
)

type Hosts struct {
	List []Host `json:"hosts"`
}

type Host struct {
	Hostname       string `json:"hostname"`
	Localhost      bool   `json:"localhost"`
	IpAddr         string `json:"ip_addr"`
	Port           string `json:"port"`
	AuthType       string `json:"auth_type"`
	User           string `json:"user"`
	Password       string `json:"password"`
	PrivateKeyFile string `json:"private_key_file"`
	Desc           string `json:"description"`
}

// CollectHosts reads hosts from given json file
func CollectHosts(hostsFile string) ([]Host, error) {
	file, err := os.ReadFile(hostsFile)
	if err != nil {
		return []Host{}, err
	}

	var hosts Hosts
	err = json.Unmarshal(file, &hosts)
	if err != nil {
		return []Host{}, err
	}

	// TODO validate hosts

	return hosts.List, nil
}

// Collect collects data from host
func (host *Host) Collect(c *collection.Collection) error {

	// TODO add file structure for host

	switch host.Localhost {
	case true:
		err := host.collectLocal(c)
		if err != nil {
			return err
		}
		break
	case false:
		err := host.collectRemote(c)
		if err != nil {
			return err
		}
	}
	return nil
}

// collect data from local
func (host *Host) collectLocal(c *collection.Collection) error {
	// TODO collect local data
	base.CollectLocal(c)
	return nil
}

// collect data from remote
func (host *Host) collectRemote(c *collection.Collection) error {
	err := host.Prepare()
	if err != nil {
		return err
	}

	client, err := host.NewConnection()
	if err != nil {
		return err
	}
	defer client.Close()
	
	// TODO replace message
	c.Log.Info("connection established")

	base.CollectRemote(client, c)

	return nil
}
