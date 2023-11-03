package collection

import (
	"github.com/mitchellh/go-ps"
	"strings"
)

// ProcessList returns array of type ps.Process. Error returns nil if no error.
func ProcessList() (processList []ps.Process, err error) {
	processList, err = ps.Processes()
	if err != nil {
		return
	}

	return
}

// ProcessListFilter returns string array of processes with given filter. Error returns nil if no error.
func ProcessListFilter(filter []string) (processList []ps.Process, err error) {
	processes, err := ProcessList()
	if err != nil {
		return
	}

	for _, process := range processes {
		for _, value := range filter {
			if strings.Contains(process.Executable(), value) {
				processList = append(processList, process)
			}
		}
	}

	return
}
