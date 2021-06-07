package util

import "os/user"

// StringInSlice matches if a string is contained in a slice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}

	return false
}

// IsPrivilegedUser returns true when the current user is root
func IsPrivilegedUser() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}

	// TODO: only works on *NIX systems
	return u.Uid == "0"
}
