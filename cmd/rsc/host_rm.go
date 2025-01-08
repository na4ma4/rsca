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

var cmdHostRemove = &cobra.Command{
	Use:   "rm <hostname> [hostname0]...[hostnameN]",
	Short: "Remove Host(s)",
	Run:   hostRemoveCommand,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	cmdHost.AddCommand(cmdHostRemove)
}

//nolint:forbidigo // Display Function
func hostRemoveCommand(_ *cobra.Command, args []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca")

	logger, _ := zapConfig().Build()
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gc := dialGRPC(ctx, cfg, logger)

	cc := api.NewAdminClient(gc)

	resp, err := cc.RemoveHost(ctx, api.RemoveHostRequest_builder{
		Names: args,
	}.Build())
	if err != nil {
		logger.Fatal("unable to send RemoveHost command to server", zap.Error(err))
	}

	fmt.Printf("Hosts Removed: %q\n", resp.GetNames())
}
