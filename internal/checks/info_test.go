package checks_test

import (
	"context"
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/checks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func generateCheck(name, cmd string) *checks.Info {
	return &checks.Info{
		Name:    name,
		Timeout: 1 * time.Second,
		Command: cmd,
		Workdir: "../..",
		Period:  10 * time.Second,
	}
}

var _ = Describe("Check", func() {
	BeforeEach(func() {
		// mainconfig.ConfigInit()
		viper.Set("general.jitter", "1s")
	})

	It("will run a successful test returning OK", func() {
		c := generateCheck("SERVICE_OK", "test/check_ok.sh")
		ctx := context.Background()
		t := time.Now()
		msg := c.Run(ctx, t)
		Expect(msg).NotTo(BeNil())
		Expect(msg.GetStatus()).To(Equal(api.Status_OK))
		Expect(msg.GetOutput()).To(Equal("Test All OK"))
		Expect(msg.GetOutputError()).To(BeEmpty())
	})

	It("will run a successful test returning WARNING", func() {
		c := generateCheck("SERVICE_WARNING", "test/check_warning.sh")
		ctx := context.Background()
		t := time.Now()
		msg := c.Run(ctx, t)
		Expect(msg).NotTo(BeNil())
		Expect(msg.GetStatus()).To(Equal(api.Status_WARNING))
		Expect(msg.GetOutput()).To(Equal("Test returns WARNING"))
		Expect(msg.GetOutputError()).NotTo(BeEmpty())
	})

	It("will run a successful test returning CRITICAL", func() {
		c := generateCheck("SERVICE_CRITICAL", "test/check_error.sh")
		ctx := context.Background()
		t := time.Now()
		msg := c.Run(ctx, t)
		Expect(msg).NotTo(BeNil())
		Expect(msg.GetStatus()).To(Equal(api.Status_CRITICAL))
		Expect(msg.GetOutput()).To(Equal("Test returns ERROR"))
		Expect(msg.GetOutputError()).NotTo(BeEmpty())
	})

	It("will run a successful test returning UNKNOWN", func() {
		c := generateCheck("SERVICE_UNKNOWN", "test/check_unknown.sh")
		ctx := context.Background()
		t := time.Now()
		msg := c.Run(ctx, t)
		Expect(msg).NotTo(BeNil())
		Expect(msg.GetStatus()).To(Equal(api.Status_UNKNOWN))
		Expect(msg.GetOutput()).To(Equal("Test returns UNKNOWN"))
		Expect(msg.GetOutputError()).NotTo(BeEmpty())
	})

	It("will timeout on long running commands", func() {
		c := generateCheck("SLEEP_TEST", "test/check_sleep.sh")
		ctx := context.Background()
		t := time.Now()
		msg := c.Run(ctx, t)
		Expect(msg).NotTo(BeNil())
		Expect(msg.GetStatus()).To(Equal(api.Status_UNKNOWN))
		Expect(msg.GetOutput()).To(ContainSubstring("Test goes to SLEEP"))
		Expect(msg.GetOutput()).NotTo(ContainSubstring("Test woke up"))
		Expect(msg.GetOutputError()).To(ContainSubstring("check timeout"))
	}, 15)

	It("will return UNKNOWN on invalid commands", func() {
		c := generateCheck("INVALID_COMMAND", "no-such-command-exists.fish")
		ctx := context.Background()
		t := time.Now()
		msg := c.Run(ctx, t)
		Expect(msg).NotTo(BeNil())
		Expect(msg.GetStatus()).To(Equal(api.Status_UNKNOWN))
		Expect(msg.GetOutputError()).To(ContainSubstring("check failed to run"))
	})
})

var _ = Describe("Check Runner", func() {
	BeforeEach(func() {
		// mainconfig.ConfigInit()
		viper.Set("general.jitter", "1ms")
		viper.Set("general.check-tick", "1ms")
		viper.Set("check.TEST.command", "test/check_ok.sh")
		viper.Set("check.TEST.period", "1ms")
	})

	It("will run checks", func() {
		var respEvent *api.EventMessage

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		cfg := config.NewViperConfigFromViper(viper.GetViper(), "rsca-not-used")
		logger := zap.NewNop()
		checkList := []*checks.Info{generateCheck("TEST", "test/check_ok.sh")}
		respChan := make(chan *api.EventMessage)

		go func() {
			for resp := range respChan {
				respEvent = resp
				cancel()
			}
		}()

		runner := checks.RunChecks(ctx, cfg, logger, checkList, respChan)
		err := runner()
		Expect(err).NotTo(HaveOccurred())

		Expect(respEvent).NotTo(BeNil())
		Expect(respEvent.GetStatus()).To(Equal(api.Status_OK))
		Expect(respEvent.GetOutput()).To(Equal("Test All OK"))
	}, 2)
})
