package common

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/na4ma4/config"
	"github.com/okzk/sdnotify"
	"go.uber.org/zap"
)

// ErrSignalReceived is returned when an interrupt or term signal is received.
var ErrSignalReceived = errors.New("signal received")

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
			case s := <-c:
				if cfg.GetBool("watchdog.enabled") {
					_ = sdnotify.Stopping()
				}

				close(c)
				cancel()

				return fmt.Errorf("%w: %s", ErrSignalReceived, s.String())
			case <-ctx.Done():
				logger.Debug("WaitForOSSignal Done()")

				if cfg.GetBool("watchdog.enabled") {
					_ = sdnotify.Stopping()

					return nil
				}

				close(c)

				return nil
			}
		}
	}
}
