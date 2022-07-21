package main

import (
	"context"
	"fmt"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//nolint:gochecknoglobals // cobra uses globals in main
var cmdHostRemove = &cobra.Command{
	Use:   "rm <hostname> [hostname0]...[hostnameN]",
	Short: "Remove Host(s)",
	Run:   hostRemoveCommand,
	Args:  cobra.MinimumNArgs(1),
}

//nolint:gochecknoinits // init is used in main for cobra
func init() {
	cmdHost.AddCommand(cmdHostRemove)
}

//nolint:forbidigo // Display Function
func hostRemoveCommand(cmd *cobra.Command, args []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca")

	logger, _ := zapConfig().Build()
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gc := dialGRPC(ctx, cfg, logger)

	cc := api.NewAdminClient(gc)

	resp, err := cc.RemoveHost(ctx, &api.RemoveHostRequest{
		Names: args,
	})
	if err != nil {
		logger.Fatal("unable to send RemoveHost command to server", zap.Error(err))
	}

	fmt.Printf("Hosts Removed: %q\n", resp.GetNames())
}
