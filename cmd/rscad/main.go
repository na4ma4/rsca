package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/na4ma4/config"
	"github.com/na4ma4/go-certprovider"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/common"
	"github.com/na4ma4/rsca/internal/mainconfig"
	"github.com/na4ma4/rsca/internal/state"
	"github.com/na4ma4/rsca/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var rootCmd = &cobra.Command{
	Use: "rscad",
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
}

func main() {
	_ = rootCmd.Execute()
}

func mainCommand(_ *cobra.Command, _ []string) {
	cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca")

	logger, _ := cfg.ZapConfig().Build()
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lis, listenErr := net.Listen("tcp", cfg.GetString("server.listen"))
	if listenErr != nil {
		logger.Fatal("failed to listen", zap.Error(listenErr))
	}

	cp, cpErr := certprovider.NewFileProvider(
		cfg.GetString("server.cert-dir"),
		certprovider.ServerProvider(),
	)
	if cpErr != nil {
		logger.Fatal("failed to get certificates", zap.Error(cpErr))
	}

	logger.Info("server listening", zap.String("bind", viper.GetString("server.listen")))

	st, stateErr := state.NewDiskState(logger, cfg.GetString("server.state-store"))
	if stateErr != nil {
		logger.Fatal("failed to create disk state storage", zap.Error(stateErr))
	}

	defer st.Close()

	// hostName := getHostname(cfg)
	eg, ctx := errgroup.WithContext(ctx)
	sapi := server.NewServer(logger, st)
	gc := grpc.NewServer(cp.ServerOption())

	api.RegisterRSCAServer(gc, sapi)
	api.RegisterAdminServer(gc, sapi)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	eg.Go(common.WaitForOSSignal(ctx, cancel, cfg, logger, c))
	eg.Go(sapi.Run(ctx, cfg))
	eg.Go(common.StateReaper(ctx, cfg, logger, st))
	eg.Go(common.ProcessWatchdog(ctx, cancel, cfg, logger))
	eg.Go(func() error { return gc.Serve(lis) })

	if cfg.GetBool("metrics.enabled") {
		go func() {
			http.Handle("/metrics", promhttp.Handler())

			srv := http.Server{
				Addr:              cfg.GetString("metrics.listen"),
				ReadTimeout:       cfg.GetDuration("metrics.timeout.read"),
				ReadHeaderTimeout: cfg.GetDuration("metrics.timeout.read-header"),
				WriteTimeout:      cfg.GetDuration("metrics.timeout.write"),
				IdleTimeout:       cfg.GetDuration("metrics.timeout.idle"),
			}

			if err := srv.ListenAndServe(); err != nil {
				logger.Debug("metrics.Listen context done", zap.Error(ctx.Err()))

				cancel()
			}
		}()
	}

	<-ctx.Done()
}

// func getHostname(cfg config.Conf) string {
// 	hostName := cfg.GetString("general.hostname")
// 	if hostName == "" {
// 		hostName, _ = os.Hostname()
// 	}

// 	return hostName
// }
