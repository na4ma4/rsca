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

// nolint: gochecknoglobals // cobra uses globals in main
var cmdTriggerAll = &cobra.Command{
	Use:     "all",
	Aliases: []string{"a"},
	Short:   "Trigger all services on a host",
	Run:     triggerAllCommand,
	Args:    cobra.NoArgs,
}

// nolint:gochecknoinits // init is used in main for cobra
func init() {
	cmdTriggerAll.PersistentFlags().StringSliceP("tags", "t", []string{},
		"tags to target, OR'd list, specified argument repeatedly to target multiple tags",
	)
	cmdTrigger.AddCommand(cmdTriggerAll)

	_ = viper.BindPFlag("trigger.all.tags", cmdTriggerAll.PersistentFlags().Lookup("tags"))
}

func triggerAllCommand(cmd *cobra.Command, args []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca")

	logger, _ := cfg.ZapConfig().Build()
	defer logger.Sync() //nolint: errcheck

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gc := dialGRPC(ctx, cfg, logger)
	cc := api.NewAdminClient(gc)

	ms := &api.Members{
		Tag: cfg.GetStringSlice("trigger.all.tags"),
	}

	r, err := cc.TriggerAll(ctx, ms)
	if err != nil {
		logger.Fatal("unable to trigger all services", zap.Error(err))
	}

	if r != nil {
		fmt.Println("Trigger message sent")

		return
	}

	fmt.Println("Trigger message failed")
}
