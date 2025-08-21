package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/na4ma4/config"
	"github.com/na4ma4/go-slogtool"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	_, logger := helpers.LogManager(slog.LevelInfo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gc := dialGRPC(ctx, cfg, logger)

	cc := api.NewAdminClient(gc)

	resp, err := cc.RemoveHost(ctx, api.RemoveHostRequest_builder{
		Names: args,
	}.Build())
	if err != nil {
		logger.ErrorContext(ctx, "unable to send RemoveHost command to server", slogtool.ErrorAttr(err))
		panic(err)
	}

	fmt.Printf("Hosts Removed: %q\n", resp.GetNames())
}
