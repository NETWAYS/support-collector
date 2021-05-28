package base

import (
	"fmt"
	"github.com/Showmax/go-fqdn"
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

	i.FQDN, err = fqdn.FqdnHostname()
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

		s[i] = uint8(chars[i])
	}

	return string(s[0:i])
}
