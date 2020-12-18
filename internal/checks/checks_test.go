package checks_test

import (
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/checks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var _ = Describe("Check Loading", func() {
	BeforeEach(func() {
		viper.Set("default.timeout", "12s")
		viper.Set("default.period", "64s")
		viper.Set("default.name-format", "uppercase")
		viper.Set("check.TEST.name", "foo")
		viper.Set("check.TEST.command", "test-command.sh")
		viper.Set("check.TEST.period", "10s")
		viper.Set("check.TEST.timeout", "11s")
		viper.Set("check.TEST.workdir", "/tmp")
		viper.Set("check.HOST.name", "PING")
		viper.Set("check.HOST.type", "host")
		viper.Set("check.HOST.command", "host-check-command.sh")
		viper.Set("check.HOST.period", "22s")
		viper.Set("check.BAD_TYPE.type", "bar")
	})

	It("will load up the TEST checks", func() {
		// checks

		cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca-not-used")
		logger := zap.NewNop()
		checkList := checks.GetChecksFromViper(cfg, viper.GetViper(), logger, "localhost.localdomain")

		Expect(checkList).To(ContainElement(&checks.Info{
			Type:     api.CheckType_SERVICE,
			Hostname: "localhost.localdomain",
			Name:     "foo",
			Command:  "test-command.sh",
			Period:   10 * time.Second,
			Timeout:  11 * time.Second,
			Workdir:  "/tmp",
		}))

		Expect(checkList).To(ContainElement(&checks.Info{
			Type:     api.CheckType_HOST,
			Hostname: "localhost.localdomain",
			Command:  "host-check-command.sh",
			Name:     "PING",
			Period:   22 * time.Second,
			Timeout:  12 * time.Second,
		}))

		Expect(checkList).To(ContainElement(&checks.Info{
			Type:     api.CheckType_SERVICE,
			Hostname: "localhost.localdomain",
			Name:     "BAD_TYPE",
			Period:   64 * time.Second,
			Timeout:  12 * time.Second,
		}))
	})
})
