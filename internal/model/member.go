package model

import (
	"time"

	"github.com/na4ma4/rsca/api"
)

type Member struct {
	ID           string        `json:"id,omitempty"`
	InternalID   string        `json:"internal_id,omitempty"`
	Name         string        `json:"name,omitempty"`
	Capability   []string      `json:"capability,omitempty"`
	Tag          []string      `json:"tag,omitempty"`
	Service      []string      `json:"service,omitempty"`
	Version      string        `json:"version,omitempty"`
	GitHash      string        `json:"git_hash,omitempty"`
	BuildDate    string        `json:"build_date,omitempty"`
	LastSeen     time.Time     `json:"last_seen,omitempty"`
	PingLatency  time.Duration `json:"ping_latency,omitempty"`
	InfoStat     *InfoStat     `json:"infostat,omitempty"`
	SystemStart  time.Time     `json:"system_start,omitempty"`
	ProcessStart time.Time     `json:"process_start,omitempty"`
	Active       bool          `json:"active,omitempty"`
	LastSeenAgo  string        `json:"lastseenago,omitempty"`
	Latency      string        `json:"latency,omitempty"`
}

func MemberFromAPI(in *api.Member) *Member {
	return &Member{
		ID:           in.GetId(),
		InternalID:   in.GetInternalId(),
		Name:         in.GetName(),
		Capability:   in.GetCapability(),
		Tag:          in.GetTag(),
		Service:      in.GetService(),
		Version:      in.GetVersion(),
		GitHash:      in.GetGitHash(),
		BuildDate:    in.GetBuildDate(),
		LastSeen:     in.GetLastSeen().AsTime(),
		PingLatency:  in.GetPingLatency().AsDuration(),
		SystemStart:  in.GetSystemStart().AsTime(),
		ProcessStart: in.GetProcessStart().AsTime(),
		Active:       in.GetActive(),
		LastSeenAgo:  in.GetLastSeenAgo(),
		Latency:      in.GetLatency(),
		InfoStat:     InfoStatFromAPI(in.GetInfoStat()),
	}
}

func (m *Member) GetName() string {
	return m.Name
}

type InfoStat struct {
	Timestamp       time.Time `json:"ts,omitempty"`
	Hostname        string    `json:"hostname,omitempty"`
	Uptime          uint64    `json:"uptime,omitempty"`
	BootTime        uint64    `json:"boottime,omitempty"`
	Procs           uint64    `json:"procs,omitempty"`
	OS              string    `json:"os,omitempty"`
	Platform        string    `json:"platform,omitempty"`
	PlatformFamily  string    `json:"platform_family,omitempty"`
	PlatformVersion string    `json:"platform_version,omitempty"`
	KernelVersion   string    `json:"kernel_version,omitempty"`
	KernelArch      string    `json:"kernel_arch,omitempty"`
	VirtSystem      string    `json:"virt_system,omitempty"`
	VirtRole        string    `json:"virt_role,omitempty"`
	HostID          string    `json:"host_id,omitempty"`
}

func InfoStatFromAPI(in *api.InfoStat) *InfoStat {
	if in == nil {
		return &InfoStat{}
	}

	return &InfoStat{
		Timestamp:       in.GetTimestamp().AsTime(),
		Hostname:        in.GetHostname(),
		Uptime:          in.GetUptime(),
		BootTime:        in.GetBootTime(),
		Procs:           in.GetProcs(),
		OS:              in.GetOs(),
		Platform:        in.GetPlatform(),
		PlatformFamily:  in.GetPlatformFamily(),
		PlatformVersion: in.GetPlatformVersion(),
		KernelVersion:   in.GetKernelVersion(),
		KernelArch:      in.GetKernelArch(),
		VirtSystem:      in.GetVirtSystem(),
		VirtRole:        in.GetVirtRole(),
		HostID:          in.GetHostId(),
	}
}
