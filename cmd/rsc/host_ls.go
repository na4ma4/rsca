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
	"github.com/na4ma4/rsca/internal/helpers"
	"github.com/na4ma4/rsca/internal/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdHostList = &cobra.Command{
	Use:   "ls",
	Short: "List Hosts",
	Run:   hostListCommand,
	Args:  cobra.NoArgs,
}

func init() {
	cmdHostList.PersistentFlags().StringP("format", "f",
		"{{.Name}}\t{{.Active}}\t{{time .LastSeen}}\t{{age .LastSeen}}\t{{.Tag}}\t{{.Capability}}\t{{age .SystemStart}}"+
			"\t{{.Service}}",
		"Output format (go template)",
	)

	_ = viper.BindPFlag("host.list.format", cmdHostList.PersistentFlags().Lookup("format"))

	cmdHost.AddCommand(cmdHostList)
}

func hostListCommand(_ *cobra.Command, _ []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca")
	logLevel := slog.LevelInfo
	if cfg.GetBool("debug") {
		logLevel = slog.LevelDebug
	}
	_, logger := helpers.LogManager(logLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gc := dialGRPC(ctx, cfg, logger)

	cc := api.NewAdminClient(gc)

	stream, err := cc.ListHosts(ctx, &api.Empty{})
	if err != nil {
		logger.ErrorContext(ctx, "unable to receive ListHosts stream from server", slogtool.ErrorAttr(err))
		panic(err)
	}

	if strings.Contains(viper.GetString("host.list.format"), "\\t") {
		viper.Set("host.list.format", strings.ReplaceAll(viper.GetString("host.list.format"), "\\t", "\t"))
	}

	if !strings.HasSuffix(viper.GetString("host.list.format"), "\n") {
		viper.Set("host.list.format", viper.GetString("host.list.format")+"\n")
	}

	tmpl, err := template.New("").Funcs(basicFunctions()).Parse(viper.GetString("host.list.format"))
	if err != nil {
		logger.ErrorContext(ctx, "unable to load template engine", slogtool.ErrorAttr(err))
		panic(err)
	}

	hostList := scrapeHostList(ctx, logger, stream)

	printHostList(ctx, logger, tmpl, false, hostList)
}

func scrapeHostList(ctx context.Context, logger *slog.Logger, stream api.Admin_ListHostsClient) []*model.Member {
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

		hostList = append(hostList, model.MemberFromAPI(in))
	}
}
