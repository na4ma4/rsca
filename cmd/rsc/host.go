package main

import (
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/na4ma4/rsca/api"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // cobra uses globals in main
var cmdHost = &cobra.Command{
	Use:     "host",
	Aliases: []string{"h"},
	Short:   "Host Commands",
}

//nolint:gochecknoinits // init is used in main for cobra
func init() {
	rootCmd.AddCommand(cmdHost)
}

func fillInAPIMember(in *api.Member) {
	if in.GetService() == nil {
		in.Service = []string{}
	} else {
		sort.Strings(in.GetService())
	}

	if in.GetCapability() == nil {
		in.Capability = []string{}
	} else {
		sort.Strings(in.GetCapability())
	}

	if in.GetTag() == nil {
		in.Tag = []string{}
	} else {
		sort.Strings(in.GetTag())
	}

	in.LastSeenAgo = time.Since(in.GetLastSeen().AsTime()).String()
	in.Latency = in.GetPingLatency().AsDuration().String()
}

//nolint:gomnd // ignore padding count.
func printHostList(tmpl *template.Template, forceHeaderAbsent bool, hostList []*api.Member) {
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
			log.Printf("error pparsing template: %s", err.Error())
		}
	}

	for _, in := range hostList {
		if err := tmpl.Execute(w, in); err != nil {
			log.Printf("error displaying host: %s", err.Error())
		}
	}

	_ = w.Flush()
}
