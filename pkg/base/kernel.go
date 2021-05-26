package base

import "syscall"

type KernelInfo struct {
	Kernel     string
	Hostname   string
	Version    string
	Release    string
	Machine    string
	Domainname string
}

func GetKernelInfo() KernelInfo {
	var uname syscall.Utsname
	syscall.Uname(&uname)

	return KernelInfo{
		Kernel:     CharsToString(uname.Sysname[:]),
		Hostname:   CharsToString(uname.Nodename[:]),
		Version:    CharsToString(uname.Version[:]),
		Release:    CharsToString(uname.Release[:]),
		Machine:    CharsToString(uname.Machine[:]),
		Domainname: CharsToString(uname.Domainname[:]),
	}
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
