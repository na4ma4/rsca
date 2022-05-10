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
func InfoWithContext(ctx context.Context, ts time.Time) (*InfoStat, error) {
	is, err := host.InfoWithContext(ctx)
	if err == nil {
		o := &InfoStat{
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

// isMatchWildcardPrefix returns true if the query string starts with a wildcard character.
func (x *Member) isMatchWildcardPrefix(query string) bool {
	return strings.HasPrefix(query, "*") || strings.HasPrefix(query, "%")
}

// isMatchWildcardPrefix returns true if the query string ends with a wildcard character.
func (x *Member) isMatchWildcardSuffix(query string) bool {
	return strings.HasSuffix(query, "*") || strings.HasSuffix(query, "%")
}

// IsMatch returns if a member matches a supplied query string
//
// Query strings can start or end in a * or % to indicate wildcard matching.
func (x *Member) IsMatch(query string) bool {
	if strings.EqualFold(query, x.Id) || strings.EqualFold(query, x.Name) {
		return true
	}

	if x.isMatchWildcardSuffix(query) {
		if strings.HasPrefix(x.Name, query[:len(query)-1]) {
			return true
		}
	}

	if x.isMatchWildcardPrefix(query) {
		if strings.HasSuffix(x.Name, query[1:]) {
			return true
		}
	}

	if x.isMatchWildcardPrefix(query) && x.isMatchWildcardSuffix(query) {
		if strings.Contains(x.Name, query[1:len(query)-1]) {
			return true
		}
	}

	return false
}
