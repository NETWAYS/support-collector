package icingaapi

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type Endpoint struct {
	Address  string `yaml:"address" json:"address"`
	Port     int    `yaml:"port" json:"port"`
	Username string `yaml:"username" json:"-"`
	Password string `yaml:"password" json:"-"`
}

// Returns new *http.Client with insecure TLS and Proxy from ENV
func newClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		Proxy:           http.ProxyFromEnvironment,
	}
	client := &http.Client{Transport: tr}

	return client
}

// IsReachable checks if the endpoint is reachable within 5 sec
func (endpoint *Endpoint) IsReachable() error {
	// try to dial tcp connection within 5 seconds
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", endpoint.Address, endpoint.Port), 5*time.Second)
	if err != nil {
		return fmt.Errorf("cant connect to endpoint '%s' within 5 seconds: %w", endpoint.Address, err)
	}
	defer conn.Close()

	return nil
}

// Request prepares a new request for the given resourcePath and executes it.
// Url for the request is build by the given resourcePath, and the Endpoint details (url => 'https://<endpoint.address>:<endpoint.port>/<resourcePath>')
//
//	A context with 10sec timeout for the request is build. BasicAuth with username and password set.
//	Returns err if something went wrong. Result is given as []byte.
func (endpoint *Endpoint) Request(resourcePath string) ([]byte, error) {
	// Return err if no username or password provided
	if endpoint.Username == "" || endpoint.Password == "" {
		return nil, fmt.Errorf("invalid or no username or password provided for api endpoint '%s'", endpoint.Address)
	}

	// Build context for the request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build with context and url
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s:%d/%s", endpoint.Address, endpoint.Port, resourcePath), nil)
	if err != nil {
		return nil, fmt.Errorf("cant build new request for '%s': %w", endpoint.Address, err)
	}

	// Set basic auth for request
	request.SetBasicAuth(endpoint.Username, endpoint.Password)

	response, err := newClient().Do(request)
	if err != nil {
		return nil, fmt.Errorf("cant make request for '%s' to '%s': %w", endpoint.Address, resourcePath, err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("cant read response for '%s' to '%s': %w", endpoint.Address, resourcePath, err)
	}

	// if response code is not '200 OK' throw error and return
	if response.Status != "200 OK" {
		return nil, fmt.Errorf("request failed with status code %s: %s", response.Status, string(body))
	}

	return body, nil
}
