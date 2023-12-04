package checks_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/checks"
	"go.uber.org/zap"
)

func TestLoadTestChecks(t *testing.T) {
	vcfg, cfg := initTestConfig()
	logger := zap.NewNop()
	checkList := checks.GetChecksFromViper(cfg, vcfg, logger, "localhost.localdomain")

	expectTestList := checks.Checks{
		&checks.Info{
			Type:     api.CheckType_SERVICE,
			Hostname: "localhost.localdomain",
			Name:     "BAD_TYPE",
			Period:   64 * time.Second,
			Timeout:  12 * time.Second,
		},
		&checks.Info{
			Type:     api.CheckType_HOST,
			Hostname: "localhost.localdomain",
			Command:  "host-check-command.sh",
			Name:     "PING",
			Period:   22 * time.Second,
			Timeout:  12 * time.Second,
		},
		&checks.Info{
			Type:     api.CheckType_SERVICE,
			Hostname: "localhost.localdomain",
			Name:     "foo",
			Command:  "test-command.sh",
			Period:   10 * time.Second,
			Timeout:  11 * time.Second,
			Workdir:  "/tmp",
		},
	}

	if diff := cmp.Diff(checkList, expectTestList, transformCheckList()); diff != "" {
		t.Errorf("checks.GetChecksFromViper: check list -got +want:\n%s", diff)
	}
}
