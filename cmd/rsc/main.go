package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/na4ma4/config"
	"github.com/na4ma4/go-certprovider"
	"github.com/na4ma4/go-slogtool"
	"github.com/na4ma4/rsca/internal/mainconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var rootCmd = &cobra.Command{
	Use: "rsc",
}

func init() {
	cobra.OnInitialize(mainconfig.ConfigInit)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug output")
	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindEnv("debug", "DEBUG")
}

func main() {
	_ = rootCmd.Execute()
}

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

func dialGRPC(ctx context.Context, cfg config.Conf, logger *slog.Logger) *grpc.ClientConn {
	serverHostName, _, _ := net.SplitHostPort(grpcServer(cfg.GetString("admin.server")))

	logger.DebugContext(ctx, "connecting to API",
		slog.String("bind", grpcServer(cfg.GetString("admin.server"))),
		slog.String("dns-name", serverHostName),
	)

	cp, err := certprovider.NewFileProvider(
		cfg.GetString("admin.cert-dir"),
		certprovider.CertProvider(),
	)
	if err != nil {
		logger.ErrorContext(ctx, "failed to get certificates", slogtool.ErrorAttr(err))
		panic(err)
	}

	gc, err := grpc.NewClient(grpcServer(cfg.GetString("admin.server")), cp.DialOption(serverHostName))
	if err != nil {
		logger.ErrorContext(ctx, "failed to connect to server", slogtool.ErrorAttr(err))
		panic(err)
	}

	return gc
}
