package common

import (
	"context"
	"time"

	"github.com/na4ma4/config"
	"github.com/okzk/sdnotify"
	"go.uber.org/zap"
)

// ProcessWatchdog sends systemd watchdog heartbeats to the systemd process.
func ProcessWatchdog(ctx context.Context, cancel context.CancelFunc, cfg config.Conf, logger *zap.Logger) func() error {
	if cfg.GetBool("watchdog.enabled") {
		logger.Info("starting watchdog")

		ticker := time.NewTicker(cfg.GetDuration("watchdog.tick"))

		return func() error {
			for {
				select {
				case <-ticker.C:
					logger.Debug("watchdog tick sent")

					if err := sdnotify.Watchdog(); err != nil {
						logger.Error("systemd watchdog returned error", zap.Error(err))
						cancel()
					}
				case <-ctx.Done():
					logger.Debug("ProcessWatchdog Done()")
					ticker.Stop()

					return ctx.Err()
				}
			}
		}
	}

	return func() error { return nil }
}
