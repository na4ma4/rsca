package server

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/na4ma4/permbits"
	"github.com/na4ma4/rsca/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ErrUnknownMessageType is returned when a message is of unknown type.
var ErrUnknownMessageType = errors.New("unknown message type")

func writeCheckResponse(logger *zap.Logger, msg *api.EventMessage) error {
	status := int32(msg.GetStatus())

	switch msg.GetType() {
	case api.CheckType_HOST:
		o := fmt.Sprintf(
			"PROCESS_HOST_CHECK_RESULT;%s;%d;%s",
			msg.GetHostname(),
			status,
			msg.GetOutput(),
		)

		return writeCommand(logger, o)
	case api.CheckType_SERVICE:
		o := fmt.Sprintf(
			"PROCESS_SERVICE_CHECK_RESULT;%s;%s;%d;%s",
			msg.GetHostname(),
			msg.GetCheck(),
			status,
			msg.GetOutput(),
		)

		return writeCommand(logger, o)
	default:
		return fmt.Errorf("%w: %d", ErrUnknownMessageType, msg.GetType())
	}
}

func writeCommand(logger *zap.Logger, command string) error {
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

	logger.Debug("writing to nagios command-file", zap.String("command", commandToWrite))

	_, err = f.WriteString(commandToWrite)

	if err != nil {
		return fmt.Errorf("write command to nagios: %w", err)
	}

	return nil
}
