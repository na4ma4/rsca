package checks_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/checks"
	"github.com/spf13/viper"
)

func TestCheckTimeout(t *testing.T) {
	t.Parallel()

	timeout := time.After(15 * time.Second)
	done := make(chan bool)

	c := generateCheck("SLEEP_TEST", "testdata/check_sleep.sh")
	ctx := context.Background()
	ts := time.Now()

	go func() {
		msg := c.Run(ctx, ts)

		expectCheckNotNil(t, msg)
		expectCheckStatus(t, msg, api.Status_UNKNOWN)
		expectCheckOutput(t, msg, "Test goes to SLEEP ...")
		expectCheckOutputError(t, msg, "check timeout")

		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("test didn't finish in time")
	case <-done:
	}
}

func TestCheckExecute(t *testing.T) {
	// shell tests don't run nicely in parallel on GitHub Actions.
	// t.Parallel()

	tests := []struct {
		name          string
		checkName     string
		executable    string
		expectStatus  api.Status
		expectMessage string
		expectError   string
	}{
		{
			"successful test return OK", "SERVICE_OK", "testdata/check_ok.sh",
			api.Status_OK, "Test All OK", "",
		},
		{
			"successful test return WARNING", "SERVICE_WARNING", "testdata/check_warning.sh",
			api.Status_WARNING, "Test returns WARNING", "exit status 1",
		},
		{
			"successful test return CRITICAL", "SERVICE_CRITICAL", "testdata/check_error.sh",
			api.Status_CRITICAL, "Test returns ERROR", "exit status 2",
		},
		{
			"successful test return UNKNOWN", "SERVICE_UNKNOWN", "testdata/check_unknown.sh",
			api.Status_UNKNOWN, "Test returns UNKNOWN", "exit status 3",
		},
		{
			"will timeout on long running commands", "SLEEP_TEST", "testdata/check_sleep.sh",
			api.Status_UNKNOWN, "Test goes to SLEEP ...", "check timeout",
		},
		{
			"will return UNKNOWN on invalid commands", "INVALID_COMMAND", "no-such-command-exists.fish",
			api.Status_UNKNOWN, "", "*not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			timeout := time.After(15 * time.Second)
			done := make(chan bool)

			go func() {
				c := generateCheck(tt.checkName, tt.executable)
				ctx := context.Background()
				ts := time.Now()

				t.Logf("Running Command: %+v", c)
				msg := c.Run(ctx, ts)

				t.Logf("Message: %+v", msg)
				expectCheckNotNil(t, msg)
				expectCheckStatus(t, msg, tt.expectStatus)
				expectCheckOutput(t, msg, tt.expectMessage)
				expectCheckOutputError(t, msg, tt.expectError)

				done <- true
			}()

			select {
			case <-timeout:
				t.Fatal("test didn't finish in time")
			case <-done:
			}
		})
	}
}

func TestCheckRun(t *testing.T) {
	t.Parallel()

	vcfg := viper.New()
	timeout := time.After(5 * time.Second)

	vcfg.Set("general.jitter", "1ms")
	vcfg.Set("general.check-tick", "1ms")
	vcfg.Set("check.TEST.command", "testdata/check_ok.sh")
	vcfg.Set("check.TEST.period", "1ms")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.NewViperConfigFromViper(vcfg, "rsca-not-used")
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, nil)
	logger := slog.New(h)
	checkList := []*checks.Info{generateCheck("TEST", "testdata/check_ok.sh")}
	respChan := make(chan *api.EventMessage)
	runner := checks.RunChecks(ctx, cfg, logger, checkList, respChan)

	go func(f func() error) {
		if err := f(); err != nil {
			t.Errorf("checks.RunChecks():runner: error, got '%s', want 'nil'", err)
		}
	}(runner)

	select {
	case respEvent := <-respChan:
		expectCheckNotNil(t, respEvent)
		expectCheckStatus(t, respEvent, api.Status_OK)
		expectCheckOutput(t, respEvent, "Test All OK")
		expectCheckOutputError(t, respEvent, "")

	case <-timeout:
		t.Fatal("test didn't finish in time")
	}
}
