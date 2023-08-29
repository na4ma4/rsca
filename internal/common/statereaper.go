package common

import (
	"context"
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/state"
	"go.uber.org/zap"
)

// StateReaper periodically checks the state store and deactivates old entries.
func StateReaper(
	ctx context.Context,
	cfg config.Conf,
	logger *zap.Logger,
	st state.State,
) func() error {
	logger.Info("starting state reaper")

	ticker := time.NewTicker(cfg.GetDuration("server.state-tick"))

	return func() error {
		for {
			select {
			case ts := <-ticker.C:
				logger.Debug("state reaper tick received")

				expireState := []string{}

				expireTime := ts.Add(-1 * cfg.GetDuration("server.state-timeout")).UTC()

				_ = st.Walk(func(in *api.Member) error {
					if in != nil && in.LastSeen != nil && in.Active && !in.LastSeen.AsTime().After(expireTime) {
						logger.Debug("adding host to inactive list",
							zap.String("rsca.client.name", in.GetName()),
							zap.Time("expireTime", expireTime),
							zap.Time("lastseen", in.LastSeen.AsTime()),
						)
						expireState = append(expireState, in.GetName())
					}

					return nil
				})

				for k := range expireState {
					logger.Info("deactivating host for inactivity", zap.String("rsca.client.name", expireState[k]))

					if err := st.DeactivateByHostname(expireState[k]); err != nil {
						logger.Error("unable to deactive member", zap.Error(err))
					}
				}
			case <-ctx.Done():
				logger.Debug("StateReaper Done()")
				ticker.Stop()

				return nil
			}
		}
	}
}
