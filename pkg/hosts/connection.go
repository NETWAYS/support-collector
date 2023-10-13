package hosts

import (
	"fmt"
	"net"
	"time"
)

// Prepare checks if the given connection details are valid and host is reachable
func (host *Host) Prepare() error {
	if !host.parseAddress() {
		return fmt.Errorf("given ip address is not valid")
	}

	err := host.testTCPConnection()
	if err != nil {
		return fmt.Errorf("%s is not reachable within 20sec: %s", host.Hostname, err)
	}

	return nil
}

// parseAddress checks if ip address is valid
func (host *Host) parseAddress() bool {
	addr := net.ParseIP(host.IpAddr)
	if addr == nil {
		return false
	}
	return true
}

// testTCPConnection checks if the host is reachable via tcp within 20 seconds
func (host *Host) testTCPConnection() error {
	_, err := net.DialTimeout("tcp", host.IpAddr+":"+host.Port, 20*time.Second)
	if err != nil {
		return err
	}

	return nil
}
