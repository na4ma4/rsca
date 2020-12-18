package checks

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Info is the details of a check.
type Info struct {
	Name     string
	Type     api.CheckType
	Hostname string
	Period   time.Duration
	NextRun  time.Time
	Command  string
	Timeout  time.Duration
	Workdir  string
}

// runCmd runs a supplied command and returns the exitcode.
func (i *Info) runCmd(cmd *exec.Cmd) (exitCode int, cmdErr error) {
	cmdErr = cmd.Run()
	if cmdErr != nil {
		// try to get the exit code
		var exitError *exec.ExitError
		if errors.As(cmdErr, &exitError) {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH, in this situation,
			// exit code could not be get, and stderr will be empty string very likely, so we use
			// the default fail code, and format err to string and set to stderr
			exitCode = 3
			cmdErr = fmt.Errorf("check failed to run: %w", cmdErr)
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	return
}

// Run executes a check and returns an api.EventMessage with the details.
func (i *Info) Run(ctx context.Context, t time.Time) *api.EventMessage {
	args := i.splitCmd()

	if i.Timeout > 0 {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, i.Timeout)
		defer cancel()
	}

	exitCode, ob, oberr, err := i.wrapCmd(ctx, args)
	status := api.ExitCodeToStatus(exitCode)
	resp := &api.EventMessage{
		Check:    i.Name,
		Hostname: i.Hostname,
		Type:     i.Type,
		Id:       uuid.New().String(),
		Output:   strings.TrimSpace(ob.String()),
		Status:   status,
	}

	switch {
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		resp.OutputError = "check timeout"
		resp.Status = api.Status_UNKNOWN
	case err != nil:
		resp.OutputError = err.Error()
	default:
		resp.OutputError = strings.TrimSpace(oberr.String())
	}

	if ts, err := ptypes.TimestampProto(t); err == nil {
		resp.RequestTimestamp = ts
	}

	if viper.GetDuration("general.jitter").Seconds() > 1 {
		// don't care about how secure the random is, it's for jitter calculations
		checkJitter := time.Duration(
			rand.Intn(int(viper.GetDuration("general.jitter").Seconds())), //nolint: gosec
		) * time.Second
		i.NextRun = time.Now().Add(i.Period).Add(checkJitter)
	} else {
		i.NextRun = time.Now().Add(i.Period)
	}

	return resp
}

// RunChecks is a routine that will cycle through the checks on a schedule and execute any pending checks.
func RunChecks(
	ctx context.Context,
	cfg config.Conf,
	logger *zap.Logger,
	checkList []*Info,
	respChan chan *api.EventMessage,
) func() error {
	ticker := time.NewTicker(cfg.GetDuration("general.check-tick"))

	return func() error {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()

				return nil
			case t := <-ticker.C:
				for i := range checkList {
					if t.After(checkList[i].NextRun) {
						respChan <- checkList[i].Run(ctx, t)
					}
				}
			}
		}
	}
}
