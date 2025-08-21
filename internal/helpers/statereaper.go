package helpers

import (
	"context"
	"log/slog"
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/go-slogtool"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/state"
)

// StateReaper periodically checks the state store and deactivates old entries.
func StateReaper(
	ctx context.Context,
	cfg config.Conf,
	logger *slog.Logger,
	st state.State,
) func() error {
	logger.InfoContext(ctx, "starting state reaper")

	ticker := time.NewTicker(cfg.GetDuration("server.state-tick"))

	return func() error {
		for {
			select {
			case ts := <-ticker.C:
				logger.DebugContext(ctx, "state reaper tick received")

				expireState := []string{}

				expireTime := ts.Add(-1 * cfg.GetDuration("server.state-timeout")).UTC()

				_ = st.Walk(func(in *api.Member) error {
					if in != nil && in.GetLastSeen() != nil &&
						in.GetActive() &&
						!in.GetLastSeen().AsTime().After(expireTime) {
						logger.Debug("adding host to inactive list",
							slog.String("rsca.client.name", in.GetName()),
							slog.Time("expireTime", expireTime),
							slog.Time("lastseen", in.GetLastSeen().AsTime()),
						)
						expireState = append(expireState, in.GetName())
					}

					return nil
				})

				for k := range expireState {
					logger.InfoContext(ctx,
						"deactivating host for inactivity",
						slog.String("rsca.client.name", expireState[k]),
					)

					if err := st.DeactivateByHostname(expireState[k]); err != nil {
						logger.ErrorContext(ctx, "unable to deactive member", slogtool.ErrorAttr(err))
					}
				}
			case <-ctx.Done():
				logger.DebugContext(ctx, "StateReaper Done()")
				ticker.Stop()

				return nil
			}
		}
	}
}
