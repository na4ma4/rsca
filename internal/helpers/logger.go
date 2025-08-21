package helpers

import (
	"io"
	"log/slog"
	"os"

	"github.com/na4ma4/go-slogtool"
	"github.com/na4ma4/go-slogtool/prettylog"
)

func LogManager(logLevel slog.Leveler) (slogtool.LogManager, *slog.Logger) {
	opts := []interface{}{
		slogtool.WithSource(false),
		slogtool.WithDefaultLevel(logLevel),
	}

	if logLevel.Level() <= slog.LevelDebug {
		opts = append(opts,
			slogtool.WithCustomHandler(func(_ string, _ io.Writer, opts *slog.HandlerOptions) slog.Handler {
				return prettylog.NewHandler(os.Stdout, opts)
			}),
		)
	}

	logmgr := slogtool.NewSlogManager(opts...)
	coreLogger := logmgr.Named("Core")

	return logmgr, coreLogger
}
