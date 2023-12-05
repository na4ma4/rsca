package checks_test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/checks"
	"github.com/spf13/viper"
)

func initTestConfig() (*viper.Viper, config.Conf) {
	vcfg := viper.New()

	vcfg.Set("default.timeout", "12s")
	vcfg.Set("default.period", "64s")
	vcfg.Set("default.name-format", "uppercase")
	vcfg.Set("check.TEST.name", "foo")
	vcfg.Set("check.TEST.command", "test-command.sh")
	vcfg.Set("check.TEST.period", "10s")
	vcfg.Set("check.TEST.timeout", "11s")
	vcfg.Set("check.TEST.workdir", "/tmp")
	vcfg.Set("check.HOST.name", "PING")
	vcfg.Set("check.HOST.type", "host")
	vcfg.Set("check.HOST.command", "host-check-command.sh")
	vcfg.Set("check.HOST.period", "22s")
	vcfg.Set("check.BAD_TYPE.type", "bar")

	cfg := config.NewViperConfigFromViper(vcfg, "rsca-not-used")
	return vcfg, cfg
}

func transformCheckList() cmp.Option {
	return cmp.Transformer("Sort", func(in checks.Checks) checks.Checks {
		out := append(checks.Checks{}, in...) // Copy input to avoid mutating it
		sort.Slice(out, func(i, j int) bool {
			return out[i].Hostname+out[i].Name < out[j].Hostname+out[j].Name
		})
		return out
	})
}

func generateCheck(name, cmd string) *checks.Info {
	return &checks.Info{
		Name:    name,
		Timeout: 5 * time.Second,
		Command: cmd,
		Workdir: "../..",
		Period:  10 * time.Second,
	}
}

func expectCheckNotNil(t *testing.T, msg *api.EventMessage) {
	t.Helper()

	if msg == nil {
		t.Error("msg, got 'nil', expect 'not nil'")
	}
}

func expectCheckStatus(t *testing.T, msg *api.EventMessage, expectStatus api.Status) {
	t.Helper()

	if v := msg.GetStatus(); v != expectStatus {
		t.Errorf(
			"msg.GetStatus(): got '%d(%s)', expect '%d(%s)'",
			v, v.String(),
			expectStatus, expectStatus.String(),
		)
	}
}

func expectCheckOutput(t *testing.T, msg *api.EventMessage, expectOutput string) {
	t.Helper()

	if trimExpect, trimOK := strings.CutPrefix(expectOutput, "*"); trimOK {
		if v := msg.GetOutput(); !strings.Contains(v, trimExpect) {
			t.Errorf(
				"msg.GetOutput(): got '%s', expected to contain '%s'",
				v,
				trimExpect,
			)
		}

		return
	}

	if v := msg.GetOutput(); v != expectOutput {
		t.Errorf(
			"msg.GetOutput(): got '%s', expect '%s'",
			v,
			expectOutput,
		)
	}
}

func expectCheckOutputError(t *testing.T, msg *api.EventMessage, expectOutputError string) {
	t.Helper()

	if trimExpect, trimOK := strings.CutPrefix(expectOutputError, "*"); trimOK {
		if v := msg.GetOutputError(); !strings.Contains(v, trimExpect) {
			t.Errorf(
				"msg.GetOutputError(): got '%s', expected to contain '%s'",
				v,
				trimExpect,
			)
		}

		return
	}

	if v := msg.GetOutputError(); v != expectOutputError {
		t.Errorf(
			"msg.GetOutputError(): got '%s', expect '%s'",
			v,
			expectOutputError,
		)
	}
}
