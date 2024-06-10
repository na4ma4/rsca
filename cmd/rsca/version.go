package main

import "github.com/dosquad/go-cliversion"

func init() {
	rootCmd.Version = cliversion.Get().VersionString()
}
