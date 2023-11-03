package connector

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
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

// NewConnection dials new ssh connection and returns *ssh.Client
func (host *Host) NewConnection() (conn *ssh.Client, err error) {
	var authentication *ssh.ClientConfig

	switch host.AuthType {
	case "password":
		authentication = buildPassAuth(host)
	case "ssh-key":
		authentication, err = buildSSHKeyAuth(host)
		if err != nil {
			return &ssh.Client{}, err
		}
		break
	default:
		return &ssh.Client{}, fmt.Errorf("invalid auth type given for host %s", host.Hostname)
	}

	conn, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", host.IpAddr, host.Port), authentication)
	if err != nil {
		return &ssh.Client{}, err
	}
	return conn, nil
}

// buildPassAuth returns ssh client config for user and password authentication
func buildPassAuth(host *Host) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(host.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

// buildSSHKeyAuth returns ssh client config for user and ssh key authentication
func buildSSHKeyAuth(host *Host) (*ssh.ClientConfig, error) {
	key, err := preparePrivateKey(host.PrivateKeyFile)
	if err != nil {
		return &ssh.ClientConfig{}, err
	}

	return &ssh.ClientConfig{
		User: host.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

// preparePrivateKey prepares the given private key fqdn for ssh authentication
func preparePrivateKey(keyfile string) (ssh.Signer, error) {
	fileContent, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, fmt.Errorf("cant read private key file")
	}

	key, err := ssh.ParsePrivateKey(fileContent)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	return key, nil
}

// NewSession opens a new session for given *ssh.Client
func (host *Host) NewSession(client *ssh.Client) (*ssh.Session, error) {
	session, err := client.NewSession()
	if err != nil {
		return &ssh.Session{}, err
	}

	return session, nil
}
