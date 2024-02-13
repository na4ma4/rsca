package main

import (
	"github.com/spf13/cobra"
)

var cmdTrigger = &cobra.Command{
	Use:     "trigger",
	Aliases: []string{"t"},
	Short:   "Trigger Commands",
}

func init() {
	rootCmd.AddCommand(cmdTrigger)
}
