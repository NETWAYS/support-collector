package main

import "fmt"

//nolint:gochecknoglobals
var (
	version = "main"
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
