package main

import (
	"fmt"
)

//nolint: gochecknoglobals // these have to be variables for the linker to change the values
var (
	version   = "dev"
	buildDate = "notset"
	gitHash   = ""
)

//nolint:gochecknoinits // init is used in main for cobra
func init() {
	rootCmd.Version = fmt.Sprintf("%s [%s] (%s)", version, gitHash, buildDate)
}
