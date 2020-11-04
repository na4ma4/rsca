package common

import (
	"context"
	"os"

	"github.com/na4ma4/config"
	"github.com/okzk/sdnotify"
	"go.uber.org/zap"
)

// WaitForOSSignal listens for a signal from the OS that the system is shutting down and sets the display to "DEAD".
func WaitForOSSignal(
	ctx context.Context,
	cancel context.CancelFunc,
	cfg config.Conf,
	logger *zap.Logger,
	c chan os.Signal,
) func() error {
	return func() error {
		for {
			select {
			case <-c:
				if cfg.GetBool("watchdog.enabled") {
					_ = sdnotify.Stopping()
				}

				close(c)
				cancel()

				return nil
			case <-ctx.Done():
				logger.Debug("WaitForOSSignal Done()")

				if cfg.GetBool("watchdog.enabled") {
					return sdnotify.Stopping()
				}

				close(c)

				return ctx.Err()
			}
		}
	}
}
