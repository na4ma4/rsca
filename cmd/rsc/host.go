package main

import (
	"context"
	"log/slog"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/na4ma4/go-slogtool"
	"github.com/na4ma4/rsca/api"
	"github.com/spf13/cobra"
)

var cmdHost = &cobra.Command{
	Use:     "host",
	Aliases: []string{"h"},
	Short:   "Host Commands",
}

func init() {
	rootCmd.AddCommand(cmdHost)
}

func fillInAPIMember(_ context.Context, in *api.Member) {
	if in.GetService() == nil {
		in.SetService([]string{})
	} else {
		sort.Strings(in.GetService())
	}

	if in.GetCapability() == nil {
		in.SetCapability([]string{})
	} else {
		sort.Strings(in.GetCapability())
	}

	if in.GetTag() == nil {
		in.SetTag([]string{})
	} else {
		sort.Strings(in.GetTag())
	}

	in.SetLastSeenAgo(time.Since(in.GetLastSeen().AsTime()).String())
	in.SetLatency(in.GetPingLatency().AsDuration().String())
}

//nolint:gomnd // ignore padding count.
func printHostList(
	ctx context.Context,
	logger *slog.Logger,
	tmpl *template.Template,
	forceHeaderAbsent bool,
	hostList []*api.Member,
) {
	sort.Slice(hostList, func(i, j int) bool { return hostList[i].GetName() < hostList[j].GetName() })

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// if !strings.Contains(viper.GetString("host.list.format"), "json") {
	if !strings.Contains(tmpl.Root.String(), "json") && !forceHeaderAbsent {
		if err := tmpl.Execute(w, map[string]interface{}{
			"Id":           "ID",
			"BuildDate":    "Build Date",
			"Capability":   "Capabilities",
			"GitHash":      "Git Hash",
			"InternalId":   "Internal ID",
			"LastSeen":     "Last Seen",
			"LastSeenAgo":  "Last Seen Ago",
			"Latency":      "Latency",
			"Name":         "Name",
			"Active":       "Active",
			"PingLatency":  "Ping Latency",
			"SystemStart":  "System Start",
			"ProcessStart": "Process Start",
			"InfoStat": map[string]string{
				"Timestamp":       "Timestamp",
				"Hostname":        "Host Name",
				"Uptime":          "Uptime",
				"BootTime":        "Boot Time",
				"Procs":           "Procs",
				"Os":              "OS",
				"Platform":        "Platform",
				"PlatformFamily":  "Platform Family",
				"PlatformVersion": "Platform Version",
				"KernelVersion":   "Kernel Version",
				"KernelArch":      "Kernel Arch",
				"VirtSystem":      "Virtual System",
				"VirtRole":        "Virtual Role",
				"HostId":          "Host ID",
			},
			"Service": "Services",
			"Tag":     "Tags",
			"Version": "Version",
		}); err != nil {
			logger.ErrorContext(ctx, "error parsing template", slogtool.ErrorAttr(err))
		}
	}

	for _, in := range hostList {
		if err := tmpl.Execute(w, in); err != nil {
			logger.ErrorContext(ctx, "error displaying host", slogtool.ErrorAttr(err))
		}
	}

	_ = w.Flush()
}
