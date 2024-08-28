package util

import "os"

// ModuleExists checks if module is installed based on given relevant paths
func ModuleExists(paths []string) bool {
	for _, path := range paths {
		_, err := os.Stat(path)
		if err == nil {
			return true
		}
	}

	return false
}
