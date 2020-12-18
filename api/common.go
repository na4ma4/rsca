package api

import (
	context "context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/shirou/gopsutil/v3/host"
)

// InfoWithContext calls shirou/gopsutil InfoWithContext and returns a native InfoStat for protobuf.
func InfoWithContext(ctx context.Context, ts time.Time) (o *InfoStat, err error) {
	is, err := host.InfoWithContext(ctx)
	if err == nil {
		o = &InfoStat{
			Hostname:        is.Hostname,
			Uptime:          is.Uptime,
			BootTime:        is.BootTime,
			Procs:           is.Procs,
			Os:              is.OS,
			Platform:        is.Platform,
			PlatformFamily:  is.PlatformFamily,
			PlatformVersion: is.PlatformVersion,
			KernelVersion:   is.KernelVersion,
			KernelArch:      is.KernelArch,
			VirtSystem:      is.VirtualizationSystem,
			VirtRole:        is.VirtualizationRole,
			HostId:          is.HostID,
		}

		if tp, err := ptypes.TimestampProto(ts); err == nil {
			o.Timestamp = tp
		}

		return o, nil
	}

	return nil, fmt.Errorf("infostat retrieval failed: %w", err)
}
