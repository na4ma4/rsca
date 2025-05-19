package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"strings"
	"text/template"

	"github.com/na4ma4/config"
	"github.com/na4ma4/go-slogtool"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/common"
	"github.com/na4ma4/rsca/internal/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdHostInfo = &cobra.Command{
	Use:   "info <id|name> [id|name]",
	Short: "Host Info",
	Run:   hostInfoCommand,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	cmdHostInfo.PersistentFlags().StringP("format", "f",
		"{{.Name}}\t{{time .LastSeen}}\t{{.LastSeenAgo}}\t{{.Tag}}\t{{.Capability}}\t{{time .SystemStart}}"+
			"\t{{time .ProcessStart}}\t{{.Service}}",
		"Output format (go template)",
	)
	cmdHostInfo.PersistentFlags().BoolP("raw", "r", false,
		"Raw output (no headers)",
	)

	_ = viper.BindPFlag("host.info.format", cmdHostInfo.PersistentFlags().Lookup("format"))
	_ = viper.BindPFlag("host.info.raw", cmdHostInfo.PersistentFlags().Lookup("raw"))

	cmdHost.AddCommand(cmdHostInfo)
}

func hostInfoCommand(_ *cobra.Command, args []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca")
	_, logger := common.LogManager(slog.LevelInfo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gc := dialGRPC(ctx, cfg, logger)

	cc := api.NewAdminClient(gc)

	stream, err := cc.ListHosts(ctx, &api.Empty{})
	if err != nil {
		logger.ErrorContext(ctx, "unable to receive ListHosts stream from server", slogtool.ErrorAttr(err))
		panic(err)
	}

	if strings.Contains(viper.GetString("host.info.format"), "\\t") {
		viper.Set("host.info.format", strings.ReplaceAll(viper.GetString("host.info.format"), "\\t", "\t"))
	}

	if !strings.HasSuffix(viper.GetString("host.info.format"), "\n") {
		viper.Set("host.info.format", viper.GetString("host.info.format")+"\n")
	}

	tmpl, err := template.New("").Funcs(basicFunctions()).Parse(viper.GetString("host.info.format"))
	if err != nil {
		logger.ErrorContext(ctx, "unable to load template engine", slogtool.ErrorAttr(err))
		panic(err)
	}

	hostList := findHostInList(ctx, logger, args, stream)

	printHostList(ctx, logger, tmpl, viper.GetBool("host.info.raw"), hostList)
}

func findHostInList(
	ctx context.Context,
	logger *slog.Logger,
	query []string,
	stream api.Admin_ListHostsClient,
) []*model.Member {
	hostList := []*model.Member{}

	for {
		in, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
				logger.DebugContext(ctx, "closing stream", slogtool.ErrorAttr(err))

				return hostList
			}

			return hostList
		}

		fillInAPIMember(ctx, in)

		for _, match := range query {
			if in.IsMatch(match) {
				hostList = append(hostList, model.MemberFromAPI(in))
			}
		}
	}
}
