package main

import (
	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals // cobra uses globals in main
var cmdHost = &cobra.Command{
	Use:     "host",
	Aliases: []string{"h"},
	Short:   "Host Commands",
}

// nolint:gochecknoinits // init is used in main for cobra
func init() {
	rootCmd.AddCommand(cmdHost)
}
