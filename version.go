package main

import "fmt"

//nolint:gochecknoglobals
var (
	version = "0.10.1"
	commit  = ""
	date    = ""
)

//goland:noinspection GoBoolExpressions
func getBuildInfo() string {
	result := version

	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}

	if date != "" {
		result = fmt.Sprintf("%s\ndate: %s", result, date)
	}

	return result
}

func getVersion() string {
	return version
}
