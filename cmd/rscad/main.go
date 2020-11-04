package main

import (
	"github.com/na4ma4/rsca/internal/mainconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//nolint: gochecknoglobals // cobra uses globals in main
var rootCmd = &cobra.Command{
	Use: "rscad",
}

//nolint:gochecknoinits // init is used in main for cobra
func init() {
	cobra.OnInitialize(mainconfig.ConfigInit)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug output")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindEnv("debug", "DEBUG")
}

func main() {
	_ = rootCmd.Execute()
}
