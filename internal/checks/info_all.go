// +build !windows

package checks

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/kballard/go-shellquote"
)

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
) (exitCode int, ob bytes.Buffer, oberr bytes.Buffer, err error) {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...) //nolint: gosec
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Dir = i.Workdir
	cmd.Stdout = &ob
	cmd.Stderr = &oberr

	if i.Timeout > 0 {
		timer := time.AfterFunc(i.Timeout, func() {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		})

		exitCode, err = i.runCmd(cmd)

		timer.Stop()
	} else {
		exitCode, err = i.runCmd(cmd)
	}

	return
}
