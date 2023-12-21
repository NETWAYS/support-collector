package icinga2

import (
	"crypto/tls"
	"fmt"
	"github.com/NETWAYS/support-collector/internal/collection"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type UserAuth struct {
	Username string
	Password string
}

// APICred saves the user and password. Provided as arguments
var APICred UserAuth

// APIEndpoints saves the FQDN or ip address for the endpoints, that will be collected. Provided as arguments.
var APIEndpoints []string

// InitAPICollection starts to collect data from the Icinga 2 API for given endpoints
func InitAPICollection(c *collection.Collection) error {
	// return if no endpoints are provided
	if len(APIEndpoints) == 0 {
		return fmt.Errorf("0 API endpoints provided. No data will be collected from remote targets")
	}
	c.Log.Info("Start collection of Icinga 2 API endpoints")

	// return if username or password is not provided
	if APICred.Username == "" || APICred.Password == "" {
		return fmt.Errorf("API Endpoints provided but username and/or password are missing")
	}

	for _, endpoint := range APIEndpoints {
		// check if endpoint is reachable
		err := endpointIsReachable(endpoint)
		if err != nil {
			c.Log.Warn(err)
			continue
		}
		c.Log.Debugf("Endpoint '%s' is reachable", endpoint)

		// collect /v1/status from endpoint
		err = collectStatus(endpoint, c)
		if err != nil {
			c.Log.Warn(err)
		}
	}

	return nil
}

// endpointIsReachable checks if the given endpoint is reachable within 5 sec
func endpointIsReachable(endpoint string) error {
	timeout := 5 * time.Second

	// try to dial tcp connection within 5 seconds
	conn, err := net.DialTimeout("tcp", endpoint, timeout)
	if err != nil {
		return fmt.Errorf("cant connect to endpoint '%s' within 5 seconds: %w", endpoint, err)
	}
	defer conn.Close()

	return nil
}

// collectStatus requests $endpoint$/v1/status with APICred and saves the json result to file
func collectStatus(endpoint string, c *collection.Collection) error {
	c.Log.Debugf("request data from endpoint '%s/v1/status'", endpoint)

	// allow insecure connections because of Icinga 2 certificates
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// build request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/v1/status", endpoint), nil)
	if err != nil {
		return fmt.Errorf("cant build new request for '%s': %w", endpoint, err)
	}

	// set authentication for request
	req.SetBasicAuth(APICred.Username, APICred.Password)

	// make request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cant requests status from '%s': %w", endpoint, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cant read from response: %w", err)
	}

	// if response code is not '200 OK' throw error and return
	if resp.Status != "200 OK" {
		return fmt.Errorf("request failed with status code %s: %s", resp.Status, string(body))
	}

	// add body to file
	c.AddFileJSON(filepath.Join(ModuleName, fmt.Sprintf("api-v1_status_%s.json", extractHostname(endpoint))), string(body))

	return nil
}

// extractsHostname takes the endpoint and extract the hostname of it
func extractHostname(endpoint string) string {
	splits := strings.Split(endpoint, ":")

	return splits[0]
}
