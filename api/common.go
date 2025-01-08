package api

import (
	context "context"
	"fmt"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// InfoWithContext calls shirou/gopsutil InfoWithContext and returns a native InfoStat for protobuf.
func InfoWithContext(ctx context.Context, ts time.Time) (*InfoStat, error) {
	is, err := host.InfoWithContext(ctx)
	if err == nil {
		o := InfoStat_builder{
			Hostname:        proto.String(is.Hostname),
			Uptime:          proto.Uint64(is.Uptime),
			BootTime:        proto.Uint64(is.BootTime),
			Procs:           proto.Uint64(is.Procs),
			Os:              proto.String(is.OS),
			Platform:        proto.String(is.Platform),
			PlatformFamily:  proto.String(is.PlatformFamily),
			PlatformVersion: proto.String(is.PlatformVersion),
			KernelVersion:   proto.String(is.KernelVersion),
			KernelArch:      proto.String(is.KernelArch),
			VirtSystem:      proto.String(is.VirtualizationSystem),
			VirtRole:        proto.String(is.VirtualizationRole),
			HostId:          proto.String(is.HostID),
			Timestamp:       timestamppb.New(ts),
		}.Build()

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
	if strings.EqualFold(query, x.GetId()) || strings.EqualFold(query, x.GetName()) {
		return true
	}

	if x.isMatchWildcardSuffix(query) {
		if strings.HasPrefix(x.GetName(), query[:len(query)-1]) {
			return true
		}
	}

	if x.isMatchWildcardPrefix(query) {
		if strings.HasSuffix(x.GetName(), query[1:]) {
			return true
		}
	}

	if x.isMatchWildcardPrefix(query) && x.isMatchWildcardSuffix(query) {
		if strings.Contains(x.GetName(), query[1:len(query)-1]) {
			return true
		}
	}

	return false
}
