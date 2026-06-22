package base

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

type KernelInfo struct {
	Kernel      string
	Hostname    string
	Version     string
	Release     string
	Machine     string
	Domainname  string
	FQDN        string
	CurrentTime string
}

func getFQDN() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	resolver := &net.Resolver{
		PreferGo: true,
	}

	addrs, err := resolver.LookupHost(ctx, hostname)
	if err != nil {
		return hostname, nil
	}

	for _, addr := range addrs {
		hosts, err := resolver.LookupAddr(ctx, addr)

		if err != nil || len(hosts) == 0 {
			continue
		}

		fqdn := strings.TrimSuffix(hosts[0], ".")

		return fqdn, nil
	}

	return hostname, nil
}

func GetKernelInfo() (i KernelInfo, err error) {
	var uname syscall.Utsname

	err = syscall.Uname(&uname)
	if err != nil {
		err = fmt.Errorf("could not load uname: %w", err)
		return
	}

	i = KernelInfo{
		Kernel:     CharsToString(uname.Sysname[:]),
		Hostname:   CharsToString(uname.Nodename[:]),
		Version:    CharsToString(uname.Version[:]),
		Release:    CharsToString(uname.Release[:]),
		Machine:    CharsToString(uname.Machine[:]),
		Domainname: CharsToString(uname.Domainname[:]),
	}

	i.FQDN, err = getFQDN()
	// return err after setting info

	i.CurrentTime = time.Now().String()

	return
}

func CharsToString(chars []int8) string {
	s := make([]byte, len(chars))

	var i int
	for ; i < len(chars); i++ {
		if chars[i] == 0 {
			break
		}

		s[i] = byte(chars[i]) //nolint: gosec
	}

	return string(s[0:i])
}
