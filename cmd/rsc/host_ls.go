package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// nolint: gochecknoglobals // cobra uses globals in main
var cmdHostList = &cobra.Command{
	Use:   "ls",
	Short: "List Hosts",
	Run:   hostListCommand,
	Args:  cobra.NoArgs,
}

// nolint:gochecknoinits // init is used in main for cobra
func init() {
	cmdHostList.PersistentFlags().StringP("format", "f",
		"{{.Name}}\t{{time .LastSeen}}\t{{.LastSeenAgo}}\t{{.Tag}}\t{{.Capability}}\t{{.Service}}",
		"Output format (go template)",
	)

	_ = viper.BindPFlag("host.list.format", cmdHostList.PersistentFlags().Lookup("format"))

	cmdHost.AddCommand(cmdHostList)
}

func hostListCommand(cmd *cobra.Command, args []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca")

	logger, _ := cfg.ZapConfig().Build()
	defer logger.Sync() //nolint: errcheck

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gc := dialGRPC(ctx, cfg, logger)

	cc := api.NewAdminClient(gc)

	stream, err := cc.ListHosts(ctx, &api.EmptyRequest{})
	if err != nil {
		logger.Fatal("unable to receive ListHosts stream from server", zap.Error(err))
	}

	if !strings.HasSuffix(viper.GetString("host.list.format"), "\n") {
		viper.Set("host.list.format", fmt.Sprintf("%s\n", viper.GetString("host.list.format")))
	}

	tmpl, err := template.New("").Funcs(basicFunctions).Parse(viper.GetString("host.list.format"))
	if err != nil {
		logger.Fatal("unable to load template engine", zap.Error(err))
	}

	for {
		in, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
				logger.Debug("closing stream", zap.Error(err))

				return
			}

			return
		}

		if in.Service == nil {
			in.Service = []string{}
		}

		if in.Capability == nil {
			in.Capability = []string{}
		}

		if in.Tag == nil {
			in.Tag = []string{}
		}

		in.LastSeenAgo = time.Since(in.LastSeen.AsTime()).String()

		_ = tmpl.Execute(os.Stdout, in)
	}
}
