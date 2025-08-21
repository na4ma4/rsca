package helpers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/na4ma4/rsca/api"
	"google.golang.org/protobuf/proto"
)

type streamServer interface {
	Send(*api.Message) error
}

// ProcessPingMessage is a common handler for processing PingMessage messages.
func ProcessPingMessage(
	ctx context.Context,
	logger *slog.Logger,
	stream streamServer,
	hostName string,
	in *api.Message,
	msg *api.PingMessage,
) error {
	logger.DebugContext(ctx, "Received PingMessage")

	if err := stream.Send(GeneratePingMessage(ctx, logger, hostName, in, msg)); err != nil {
		return fmt.Errorf("unable to send ping message: %w", err)
	}

	return nil
}

// GeneratePingMessage is a common handler for generating PingMessage messages.
func GeneratePingMessage(
	ctx context.Context,
	logger *slog.Logger,
	hostName string,
	in *api.Message,
	msg *api.PingMessage,
) *api.Message {
	logger.DebugContext(ctx, "Received PingMessage")

	return api.Message_builder{
		Envelope: api.Envelope_builder{
			Sender: api.Member_builder{
				Name: proto.String(hostName),
			}.Build(),
			Recipient: api.RecipientBySender(in.GetEnvelope().GetSender()),
		}.Build(),
		PongMessage: api.PongMessage_builder{
			Id:       proto.String(msg.GetId()),
			StreamId: proto.String(msg.GetStreamId()),
			Ts:       msg.GetTs(),
		}.Build(),
	}.Build()
}
