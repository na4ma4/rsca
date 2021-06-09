package api

import (
	context "context"
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			Timestamp:       timestamppb.New(ts),
		}

		return o, nil
	}

	return nil, fmt.Errorf("infostat retrieval failed: %w", err)
}

func (x *Member) IsMatch(query string) bool {
	if strings.EqualFold(query, x.Id) || strings.EqualFold(query, x.Name) {
		return true
	}

	if strings.HasSuffix(query, "*") || strings.HasSuffix(query, "%") {
		if strings.HasPrefix(x.Name, query[:len(query)-1]) {
			return true
		}
	}

	if strings.HasPrefix(query, "*") || strings.HasPrefix(query, "%") {
		if strings.HasSuffix(x.Name, query[1:]) {
			return true
		}
	}

	return false
}
