package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/na4ma4/go-permbits"
	"github.com/na4ma4/rsca/api"
	"github.com/spf13/viper"
)

// ErrUnknownMessageType is returned when a message is of unknown type.
var ErrUnknownMessageType = errors.New("unknown message type")

func writeCheckResponse(ctx context.Context, logger *slog.Logger, msg *api.EventMessage) error {
	status := int32(msg.GetStatus())

	switch msg.GetType() {
	case api.CheckType_HOST:
		o := fmt.Sprintf(
			"PROCESS_HOST_CHECK_RESULT;%s;%d;%s",
			msg.GetHostname(),
			status,
			msg.GetOutput(),
		)

		return writeCommand(ctx, logger, o)
	case api.CheckType_SERVICE:
		o := fmt.Sprintf(
			"PROCESS_SERVICE_CHECK_RESULT;%s;%s;%d;%s",
			msg.GetHostname(),
			msg.GetCheck(),
			status,
			msg.GetOutput(),
		)

		return writeCommand(ctx, logger, o)
	default:
		return fmt.Errorf("%w: %d", ErrUnknownMessageType, msg.GetType())
	}
}

func writeCommand(ctx context.Context, logger *slog.Logger, command string) error {
	command = strings.TrimSpace(command)
	commandToWrite := fmt.Sprintf("[%d] %s\n", time.Now().Unix(), command)

	f, err := os.OpenFile(
		viper.GetString("nagios.command-file"),
		os.O_APPEND|os.O_WRONLY,
		permbits.UserRead+permbits.UserWrite,
	)
	if err != nil {
		return fmt.Errorf("open command file for nagios: %w", err)
	}

	defer func() { _ = f.Close() }()

	logger.DebugContext(ctx, "writing to nagios command-file", slog.String("command", commandToWrite))

	_, err = f.WriteString(commandToWrite)
	if err != nil {
		return fmt.Errorf("write command to nagios: %w", err)
	}

	return nil
}
