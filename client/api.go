package client

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/checks"
	"github.com/na4ma4/rsca/internal/common"
	"github.com/na4ma4/rsca/internal/register"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Client is a api.RSCAClient for co-ordinating requests from the server.
type Client struct {
	logger   *zap.Logger
	hostname string
	checks   checks.Checks
	inbox    chan *api.Message
	outbox   chan *api.Message
	register *register.Message
}

// NewClient returns a setup api.RSCAClient.
func NewClient(logger *zap.Logger, hostName string, checkList checks.Checks) *Client {
	return &Client{
		logger:   logger,
		hostname: hostName,
		checks:   checkList,
		inbox:    make(chan *api.Message),
		outbox:   make(chan *api.Message),
	}
}

func (c *Client) streamMessages(cancel context.CancelFunc, stream api.RSCA_PipeClient) {
	for {
		in, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			c.logger.Debug("EOF found, closing channel")
			cancel()

			return
		} else if s, ok := status.FromError(err); err != nil && ok {
			if s.Code() == codes.Unavailable {
				c.logger.Warn("server has gone away", zap.Error(err))
				cancel()

				return
			}
		} else if err != nil {
			c.logger.Error("failed to receive a note", zap.Reflect("msg", in), zap.Error(err))
			cancel()

			return
		}

		c.inbox <- in
	}
}

// Pipe processes the stream and the outbox back to the server.
func (c *Client) Pipe(
	ctx context.Context,
	cancel context.CancelFunc,
	cfg config.Conf,
	stream api.RSCA_PipeClient,
) func() error {
	registrationTicker := time.NewTicker(cfg.GetDuration("general.registration-interval"))

	go c.streamMessages(cancel, stream)

	return func() error {
		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("context cancelled")
				registrationTicker.Stop()
				close(c.outbox)
				close(c.inbox)

				return nil
			case _, ok := <-registrationTicker.C:
				if ok {
					go c.SendRepeatRegistration(ctx)
				}
			case out, ok := <-c.outbox:
				if ok {
					if err := stream.Send(out); err != nil {
						c.logger.Error("unable to send message", zap.Error(err))
					}
				}
			}
		}
	}
}

// processUpdateAll processes a trigger all message.
func (c *Client) processUpdateAll() {
	c.logger.Debug("processUpdateAll() called")
	c.checks.NextRun(time.Time{})
}

// processRepeatRegister processes a repeat-registration request message.
func (c *Client) processRepeatRegister(ctx context.Context) {
	c.logger.Debug("processRepeatRegister() called")
	c.SendRepeatRegistration(ctx)
}

// SendRepeatRegistration sends the registration message to the server.
func (c *Client) SendRepeatRegistration(ctx context.Context) {
	c.logger.Debug("sending repeat registration message")
	c.register.UpdateInfoStat(ctx)

	c.outbox <- &api.Message{
		Envelope: &api.Envelope{Sender: c.register.Member(), Recipient: api.MembersByID("_server")},
		Message:  &api.Message_MemberUpdateMessage{MemberUpdateMessage: c.register.UpdateMessage()},
	}
}

// // Send adds a message to the outbox to be sent, may block if channel is full.
// func (c *Client) Send(msg *api.Message) {
// 	c.inbox <- msg
// }

func (c *Client) wrapEventMessage(in *api.EventMessage) *api.Message {
	return &api.Message{
		Envelope: &api.Envelope{
			Recipient: api.MembersByID("_server"),
			Sender:    c.register.Member(),
		},
		Message: &api.Message_EventMessage{
			EventMessage: in,
		},
	}
}

// RunEvents runs as a go routine that processes the response channel and creates messages to add to the outbox.
//nolint:cyclop // don't see a way to make this much more simpler without making it less readable.
func (c *Client) RunEvents(
	ctx context.Context,
	regmsg *register.Message,
	respChan chan *api.EventMessage,
) func() error {
	c.register = regmsg

	return func() error {
		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("context cancelled")

				return nil
			case in, ok := <-respChan:
				if !ok {
					return common.ErrChannelClosed
				}

				c.outbox <- c.wrapEventMessage(in)
			case in, ok := <-c.inbox:
				if !ok {
					return common.ErrChannelClosed
				}

				switch msg := in.Message.(type) {
				case *api.Message_PingMessage:
					go func() {
						c.outbox <- common.GeneratePingMessage(c.logger, c.hostname, in, msg)
					}()
				case *api.Message_TriggerAllMessage:
					go c.processUpdateAll()
				case *api.Message_RepeatRegistrationMessage:
					go c.processRepeatRegister(ctx)
				default:
					c.logger.Info("Received unhandled message", zap.Reflect("message", in))
				}
				c.logger.Debug("message processing finished")
			}
		}
	}
}
