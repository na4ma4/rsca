package checks

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/kballard/go-shellquote"
	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
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
//nolint:nestif // it might be "deeply nested", but it's readable and confines this code to this function.
func (i *Info) runCmd(cmd *exec.Cmd) (exitCode int, cmdErr error) {
	cmdErr = cmd.Run()
	if cmdErr != nil {
		// try to get the exit code
		var exitError *exec.ExitError
		if errors.As(cmdErr, &exitError) {
			if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = ws.ExitStatus()
			}
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH, in this situation,
			// exit code could not be get, and stderr will be empty string very likely, so we use
			// the default fail code, and format err to string and set to stderr
			exitCode = 3
			cmdErr = fmt.Errorf("check failed to run: %w", cmdErr)
		}
	} else if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
		// success, exitCode should be 0 if go is ok
		exitCode = ws.ExitStatus()
	}

	return
}

// splitCmd uses `shellquote` on non windows platforms.
func (i *Info) splitCmd() (o []string) {
	o, err := shellquote.Split(i.Command)
	if err != nil {
		o = strings.Split(i.Command, " ")
	}

	return
}

// wrapCmd [!windows] uses syscall.Kill to kill process group for check.
func (i *Info) wrapCmd(
	ctx context.Context,
	args []string,
) (exitCode int, ob *bytes.Buffer, oberr *bytes.Buffer, err error) {
	wg := sync.WaitGroup{}
	ob = bytes.NewBuffer(nil)
	oberr = bytes.NewBuffer(nil)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...) //nolint: gosec
	cmd.Dir = i.Workdir

	if pb, err := cmd.StdoutPipe(); err == nil {
		i.ioCopyWaitGroup(&wg, ob, pb)
	}

	if pberr, err := cmd.StderrPipe(); err == nil {
		i.ioCopyWaitGroup(&wg, oberr, pberr)
	}

	exitCode, err = i.runCmd(cmd)

	wg.Wait()

	return exitCode, ob, oberr, err
}

// ioCopyWaitGroup adds one worker to a waitgroup and runs an io.Copy until completed,
// once completed it will call waitgroup.Done().
func (i *Info) ioCopyWaitGroup(wg *sync.WaitGroup, dst io.Writer, src io.Reader) {
	wg.Add(1)

	go func() {
		_, _ = io.Copy(dst, src)

		wg.Done()
	}()
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
		Check:            i.Name,
		Hostname:         i.Hostname,
		Type:             i.Type,
		Id:               uuid.New().String(),
		Output:           strings.TrimSpace(ob.String()),
		Status:           status,
		RequestTimestamp: timestamppb.New(t),
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
