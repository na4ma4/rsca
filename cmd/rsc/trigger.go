package main

import (
	"github.com/spf13/cobra"
)

// nolint: gochecknoglobals // cobra uses globals in main
var cmdTrigger = &cobra.Command{
	Use:     "trigger",
	Aliases: []string{"t"},
	Short:   "Trigger Commands",
}

// nolint:gochecknoinits // init is used in main for cobra
func init() {
	rootCmd.AddCommand(cmdTrigger)
}
