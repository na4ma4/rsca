package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/client"
	"github.com/na4ma4/rsca/internal/certs"
	"github.com/na4ma4/rsca/internal/checks"
	"github.com/na4ma4/rsca/internal/common"
	"github.com/na4ma4/rsca/internal/mainconfig"
	"github.com/na4ma4/rsca/internal/register"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

//nolint:gochecknoglobals // cobra uses globals in main
var rootCmd = &cobra.Command{
	Use: "rsca",
	Run: mainCommand,
}

//nolint:gochecknoinits // init is used in main for cobra
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

func mainCommand(cmd *cobra.Command, args []string) {
	cfg := config.NewViperConfDFromViper(viper.GetViper(), "/etc/nagios/rsca.d/", "rsca")

	logger, _ := cfg.ZapConfig().Build()
	defer logger.Sync() //nolint:errcheck

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	serverHostName, _, _ := net.SplitHostPort(grpcServer(cfg.GetString("client.server")))

	if cfg.GetString("client.sni") != "" {
		serverHostName = cfg.GetString("client.sni")
	}

	logger.Debug("Connecting to API", zap.String("bind", grpcServer(cfg.GetString("client.server"))),
		zap.String("dns-name", serverHostName))

	cp, err := certs.NewFileCertificateProvider(
		cfg.GetString("client.cert-dir"),
		certs.CertProviderFromString(cfg.GetString("client.cert-type")),
	)
	checkErrFatal(err, logger, "failed to get certificates")

	gc, err := grpc.DialContext(ctx, grpcServer(cfg.GetString("client.server")), cp.DialOption(serverHostName))
	checkErrFatal(err, logger, "failed to connect to server")

	rc := api.NewRSCAClient(gc)
	respChan := make(chan *api.EventMessage)

	stream, err := rc.Pipe(ctx)
	checkErrFatal(err, logger, "unable to create stream")

	hostName := getHostname(cfg)
	checkList := checks.GetChecksFromViper(cfg, viper.GetViper(), logger, hostName)
	cl := client.NewClient(logger, hostName, checkList)
	regmsg := register.New(cfg, hostName, version, date, commit, checkList, time.Now())
	streamMsg := &api.Message{
		Envelope: &api.Envelope{Sender: regmsg.Member(), Recipient: api.MembersByID("_server")},
		Message:  &api.Message_RegisterMessage{RegisterMessage: regmsg.Message()},
	}
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	eg.Go(cl.Pipe(ctx, cancel, cfg, stream))
	eg.Go(checks.RunChecks(ctx, cfg, logger, checkList, respChan))
	eg.Go(common.WaitForOSSignal(ctx, cancel, cfg, logger, c))
	eg.Go(common.ProcessWatchdog(ctx, cancel, cfg, logger))
	eg.Go(cl.RunEvents(ctx, regmsg, respChan))

	if err = stream.Send(streamMsg); err != nil {
		logger.Fatal("unable to register with server", zap.Error(err))
	}

	if err := eg.Wait(); err != nil {
		logger.Fatal("routine returned error", zap.Error(err))
	}
}

func checkErrFatal(err error, logger *zap.Logger, msg string) {
	if err != nil {
		logger.Fatal(msg, zap.Error(err))
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
		return fmt.Sprintf("%s:5888", server)
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
