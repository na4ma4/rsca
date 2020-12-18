// +build windows

package checks

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

// splitCmd can not use `shellquote` windows, so falls back to splitting on spaces.
func (i *Info) splitCmd() (o []string) {
	return strings.Split(i.Command, " ")
}

// wrapCmd [windows] relies on CommandContext to handle timing out processes, this is untested.
func (i *Info) wrapCmd(ctx context.Context, args []string) (exitCode int, ob bytes.Buffer, oberr bytes.Buffer, err error) {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...) //nolint: gosec
	cmd.Dir = i.Workdir
	cmd.Stdout = &ob
	cmd.Stderr = &oberr

	exitCode, err = i.runCmd(cmd)

	return
}
