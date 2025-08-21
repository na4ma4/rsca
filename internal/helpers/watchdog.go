package helpers

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/go-slogtool"
	"github.com/okzk/sdnotify"
)

// ErrWatchdogFailed is returned when the systemd watchdog tick fails.
var ErrWatchdogFailed = errors.New("systemd watchdog failed")

// ProcessWatchdog sends systemd watchdog heartbeats to the systemd process.
func ProcessWatchdog(
	ctx context.Context,
	cancel context.CancelFunc,
	cfg config.Conf,
	logger *slog.Logger,
) func() error {
	if cfg.GetBool("watchdog.enabled") {
		logger.InfoContext(ctx, "starting watchdog")

		ticker := time.NewTicker(cfg.GetDuration("watchdog.tick"))

		return func() error {
			for {
				select {
				case <-ticker.C:
					logger.DebugContext(ctx, "watchdog tick sent")

					if err := sdnotify.Watchdog(); err != nil {
						logger.ErrorContext(ctx, "systemd watchdog returned error", slogtool.ErrorAttr(err))
						cancel()

						return ErrWatchdogFailed
					}
				case <-ctx.Done():
					logger.DebugContext(ctx, "ProcessWatchdog Done()")
					ticker.Stop()

					return nil
				}
			}
		}
	}

	return func() error { return nil }
}
