package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/client"
	"github.com/na4ma4/rsca/internal/certs"
	"github.com/na4ma4/rsca/internal/checks"
	"github.com/na4ma4/rsca/internal/common"
	"github.com/na4ma4/rsca/internal/mainconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

//nolint: gochecknoglobals // cobra uses globals in main
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
}

func main() {
	_ = rootCmd.Execute()
}

func mainCommand(cmd *cobra.Command, args []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca")

	logger, _ := cfg.ZapConfig().Build()
	defer logger.Sync() //nolint: errcheck

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	serverHostName, _, _ := net.SplitHostPort(grpcServer(cfg.GetString("client.server")))

	logger.Debug("Connecting to API",
		zap.String("bind", grpcServer(cfg.GetString("client.server"))),
		zap.String("dns-name", serverHostName),
	)

	cp, err := certs.NewFileCertificateProvider(cfg.GetString("client.cert-dir"), false)
	if err != nil {
		logger.Fatal("failed to get certificates", zap.Error(err))
	}

	gc, err := grpc.DialContext(ctx, grpcServer(cfg.GetString("client.server")), cp.DialOption(serverHostName))
	if err != nil {
		logger.Fatal("failed to connect to server", zap.Error(err))
	}

	logger.Debug("server hostname", zap.String("server.host-name", serverHostName))

	rc := api.NewRSCAClient(gc)
	respChan := make(chan *api.EventMessage)

	stream, err := rc.Pipe(ctx)
	if err != nil {
		logger.Fatal("unable to create stream", zap.Error(err))
	}

	hostName := getHostname(cfg)
	checkList := checks.GetChecksFromViper(cfg, viper.GetViper(), logger, hostName)
	cl := client.NewClient(logger, hostName, checkList)
	regmsg := registerMsg(cfg, hostName, checkList)
	ms := &api.Member{Name: hostName}
	streamMsg := &api.Message{
		Envelope: &api.Envelope{Sender: ms, Recipient: api.MembersByID("_server")},
		Message:  &api.Message_RegisterMessage{RegisterMessage: regmsg},
	}
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	eg.Go(cl.Pipe(ctx, cancel, stream))
	eg.Go(checks.RunChecks(ctx, cfg, logger, checkList, respChan))
	eg.Go(common.WaitForOSSignal(ctx, cancel, cfg, logger, c))
	eg.Go(cl.RunEvents(ctx, ms, respChan))

	if err = stream.Send(streamMsg); err != nil {
		logger.Fatal("unable to register with server", zap.Error(err))
	}

	if err := eg.Wait(); err != nil {
		logger.Error("routine returned error", zap.Error(err))
	}
}

func registerMsg(cfg config.Conf, hostName string, checkList checks.Checks) *api.RegisterMessage {
	checkNames := []string{}

	for _, check := range checkList {
		if check.Type == api.CheckType_SERVICE {
			checkNames = append(checkNames, check.Name)
		}
	}

	return &api.RegisterMessage{
		Member: &api.Member{
			Id:         uuid.New().String(),
			Name:       hostName,
			Capability: []string{"client"},
			Service:    checkNames,
			Tag:        cfg.GetStringSlice("general.tags"),
		},
	}
}

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
