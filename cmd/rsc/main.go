package main

import (
	"context"
	"fmt"
	"net"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/internal/certs"
	"github.com/na4ma4/rsca/internal/mainconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

//nolint:gochecknoglobals // cobra uses globals in main
var rootCmd = &cobra.Command{
	Use: "rsc",
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

func zapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func zapConfig() zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zapEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func dialGRPC(ctx context.Context, cfg config.Conf, logger *zap.Logger) *grpc.ClientConn {
	serverHostName, _, _ := net.SplitHostPort(grpcServer(cfg.GetString("admin.server")))

	logger.Debug("connecting to API",
		zap.String("bind", grpcServer(cfg.GetString("admin.server"))),
		zap.String("dns-name", serverHostName),
	)

	cp, err := certs.NewFileCertificateProvider(
		cfg.GetString("admin.cert-dir"),
		certs.CertProviderFromString(cfg.GetString("admin.cert-type")),
	)
	if err != nil {
		logger.Fatal("failed to get certificates", zap.Error(err))
	}

	gc, err := grpc.DialContext(ctx, grpcServer(cfg.GetString("admin.server")), cp.DialOption(serverHostName))
	if err != nil {
		logger.Fatal("failed to connect to server", zap.Error(err))
	}

	return gc
}
