package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dosquad/go-cliversion"
	"github.com/na4ma4/config"
	"github.com/na4ma4/go-certprovider"
	"github.com/na4ma4/go-slogtool"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/client"
	"github.com/na4ma4/rsca/internal/checks"
	"github.com/na4ma4/rsca/internal/common"
	"github.com/na4ma4/rsca/internal/mainconfig"
	"github.com/na4ma4/rsca/internal/register"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var rootCmd = &cobra.Command{
	Use: "rsca",
	Run: mainCommand,
}

func init() {
	cobra.OnInitialize(mainconfig.ConfigInit)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug output")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindEnv("debug", "DEBUG")

	rootCmd.PersistentFlags().Bool("watchdog", false, "Enable systemd watchdog functionality")
	_ = viper.BindPFlag("watchdog.enabled", rootCmd.PersistentFlags().Lookup("watchdog"))
	_ = viper.BindEnv("watchdog.enabled", "WATCHDOG")

	rootCmd.PersistentFlags().String("config-path", "/etc/nagios/rsca.d", "Configuration path to use for config files")
	_ = viper.BindPFlag("config.path", rootCmd.PersistentFlags().Lookup("config-path"))
	_ = viper.BindEnv("config.path", "CONFIG_PATH")
}

func main() {
	_ = rootCmd.Execute()
}

func mainCommand(_ *cobra.Command, _ []string) {
	cfg := config.NewViperConfDFromViper(viper.GetViper(), "/etc/nagios/rsca.d/", "rsca")
	_, logger := common.LogManager(slog.LevelInfo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	serverHostName, _, _ := net.SplitHostPort(grpcServer(cfg.GetString("client.server")))

	if cfg.GetString("client.sni") != "" {
		serverHostName = cfg.GetString("client.sni")
	}

	logger.DebugContext(ctx, "Connecting to API", slog.String("bind", grpcServer(cfg.GetString("client.server"))),
		slog.String("dns-name", serverHostName))

	cp, cpErr := certprovider.NewFileProvider(
		cfg.GetString("client.cert-dir"),
		certprovider.ProviderFromString(cfg.GetString("client.cert-type"), certprovider.ClientProvider()),
	)
	checkErrFatal(cpErr, logger, "failed to get certificates")

	gc, gcErr := grpc.NewClient(grpcServer(cfg.GetString("client.server")), cp.DialOption(serverHostName))
	checkErrFatal(gcErr, logger, "failed to connect to server")

	rc := api.NewRSCAClient(gc)
	respChan := make(chan *api.EventMessage)

	stream, streamErr := rc.Pipe(ctx)
	checkErrFatal(streamErr, logger, "unable to create stream")

	hostName := getHostname(cfg)
	checkList := checks.GetChecksFromViper(cfg, viper.GetViper(), logger, hostName)
	cl := client.NewClient(logger, hostName, checkList)
	regmsg := register.New(cfg, hostName, cliversion.Get(), checkList, time.Now())
	streamMsg := api.Message_builder{
		Envelope:        api.Envelope_builder{Sender: regmsg.Member(), Recipient: api.MembersByID("_server")}.Build(),
		RegisterMessage: regmsg.Message(),
	}.Build()
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	eg.Go(cl.Pipe(ctx, cancel, cfg, stream))
	eg.Go(checks.RunChecks(ctx, cfg, logger, checkList, respChan))
	eg.Go(common.WaitForOSSignal(ctx, cancel, cfg, logger, c))
	eg.Go(common.ProcessWatchdog(ctx, cancel, cfg, logger))
	eg.Go(cl.RunEvents(ctx, cancel, regmsg, respChan))

	if err := stream.Send(streamMsg); err != nil {
		logger.ErrorContext(ctx, "unable to register with server", slogtool.ErrorAttr(err))
	}

	if err := eg.Wait(); err != nil {
		logger.ErrorContext(ctx, "routine returned error", slogtool.ErrorAttr(err))
	}
}

func checkErrFatal(err error, logger *slog.Logger, msg string) {
	if err != nil {
		logger.Error(msg, slogtool.ErrorAttr(err))
	}
}

// Replaced with register package
// func registerMsg(
// 	cfg config.Conf,
// 	hostName string,
// 	checkList checks.Checks,
// 	startTime time.Time,
// ) *api.RegisterMessage {
// 	checkNames := []string{}

// 	for _, check := range checkList {
// 		if check.Type == api.CheckType_SERVICE {
// 			checkNames = append(checkNames, check.Name)
// 		}
// 	}

// 	mb := &api.Member{
// 		Id:         uuid.New().String(),
// 		Name:       hostName,
// 		Capability: []string{"client", fmt.Sprintf("rsca-%s", version)},
// 		Service:    checkNames,
// 		Tag:        cfg.GetStringSlice("general.tags"),
// 		Version:    version,
// 		BuildDate:  buildDate,
// 		GitHash:    gitHash,
// 	}

// 	if ts, err := ptypes.TimestampProto(startTime); err == nil {
// 		mb.ProcessStart = ts
// 	}

// 	if ut, err := host.BootTimeWithContext(context.Background()); err == nil {
// 		if ts, err := ptypes.TimestampProto(time.Unix(int64(ut), 0)); err == nil {
// 			mb.SystemStart = ts
// 		}
// 	}

// 	if is, err := api.InfoWithContext(context.Background(), time.Now()); err == nil {
// 		mb.InfoStat = is
// 	}

// 	return &api.RegisterMessage{
// 		Member: mb,
// 	}
// }

func grpcServer(server string) string {
	host, port, err := net.SplitHostPort(server)
	if err != nil {
		return server + ":5888"
	}

	if port == "" {
		port = "5888"
	}

	return fmt.Sprintf("%s:%s", host, port)
}

func getHostname(cfg config.Conf) string {
	hostName := cfg.GetString("general.hostname")
	if hostName == "" {
		hostName, _ = os.Hostname()
	}

	return hostName
}
