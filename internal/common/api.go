package common

import (
	"github.com/na4ma4/rsca/api"
	"go.uber.org/zap"
)

type streamServer interface {
	Send(*api.Message) error
}

// ProcessPingMessage is a common handler for processing PingMessage messages.
func ProcessPingMessage(logger *zap.Logger, stream streamServer, in *api.Message, msg *api.Message_PingMessage) error {
	logger.Debug("Received PingMessage")

	return stream.Send(&api.Message{
		Envelope: &api.Envelope{
			Sender:    &api.Member{},
			Recipient: api.RecipientBySender(in.Envelope.Sender),
		},
		Message: &api.Message_PongMessage{
			PongMessage: &api.PongMessage{
				Id: msg.PingMessage.GetId(),
				Ts: msg.PingMessage.GetTs(),
			},
		},
	})
}
