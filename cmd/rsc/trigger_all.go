package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cmdTriggerAll = &cobra.Command{
	Use:     "all [options ...] [host...] [hostN]",
	Aliases: []string{"a"},
	Short:   "Trigger all services on a host",
	Run:     triggerAllCommand,
	Args:    cobra.MinimumNArgs(0),
}

func init() {
	cmdTrigger.AddCommand(cmdTriggerAll)
	cmdTriggerAll.PersistentFlags().BoolP("info", "i", false,
		"request system infostat update instead of services",
	)
	cmdTriggerAll.PersistentFlags().StringSliceP("tags", "t", []string{},
		"tags to target, OR'd list, specified argument repeatedly to target multiple tags",
	)
	cmdTriggerAll.PersistentFlags().StringSliceP("services", "s", []string{},
		"services to target, OR'd list, specified argument repeatedly to target multiple services",
	)
	cmdTriggerAll.PersistentFlags().StringSliceP("capabilities", "c", []string{},
		"capabilities to target, OR'd list, specified argument repeatedly to target multiple capabilities",
	)

	_ = viper.BindPFlag("trigger.all.info", cmdTriggerAll.PersistentFlags().Lookup("info"))
	_ = viper.BindPFlag("trigger.all.tags", cmdTriggerAll.PersistentFlags().Lookup("tags"))
	_ = viper.BindPFlag("trigger.all.services", cmdTriggerAll.PersistentFlags().Lookup("services"))
	_ = viper.BindPFlag("trigger.all.capabilities", cmdTriggerAll.PersistentFlags().Lookup("capabilities"))
}

var errTriggerFailed = errors.New("trigger failed")

//nolint:forbidigo // Display Function
func triggerAllCommand(_ *cobra.Command, args []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca")

	logger, _ := zapConfig().Build()
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gc := dialGRPC(ctx, cfg, logger)
	cc := api.NewAdminClient(gc)

	ms := api.Members_builder{
		Tag:        cfg.GetStringSlice("trigger.all.tags"),
		Service:    cfg.GetStringSlice("trigger.all.services"),
		Capability: cfg.GetStringSlice("trigger.all.capabilities"),
	}.Build()

	if len(args) > 0 {
		ms.SetName(args)
	}

	if cfg.GetBool("trigger.all.info") { //nolint:nestif // removing nesting harms readability.
		r, reqErr := cc.TriggerInfo(ctx, ms)
		if reqErr != nil {
			logger.Fatal("unable to trigger all services", zap.Error(reqErr))
		}

		if err := triggerAllCommandInfo(r); err == nil {
			return
		}
	} else {
		r, reqErr := cc.TriggerAll(ctx, ms)
		if reqErr != nil {
			logger.Fatal("unable to trigger all services", zap.Error(reqErr))
		}

		if err := triggerAllCommandServices(r); err == nil {
			return
		}
	}

	fmt.Println("Trigger message failed")
}

//nolint:forbidigo // Display Function
func triggerAllCommandInfo(r *api.TriggerInfoResponse) error {
	if r != nil {
		fmt.Printf("Trigger message sent to %d hosts\n", len(r.GetNames()))

		for _, h := range r.GetNames() {
			fmt.Println(h)
		}

		return nil
	}

	return errTriggerFailed
}

//nolint:forbidigo // Display Function
func triggerAllCommandServices(r *api.TriggerAllResponse) error {
	if r != nil {
		fmt.Printf("Trigger message sent to %d hosts\n", len(r.GetNames()))

		for _, h := range r.GetNames() {
			fmt.Println(h)
		}

		return nil
	}

	return errTriggerFailed
}
