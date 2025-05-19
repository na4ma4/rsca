package client

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"time"

	"github.com/na4ma4/config"
	"github.com/na4ma4/go-slogtool"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/checks"
	"github.com/na4ma4/rsca/internal/common"
	"github.com/na4ma4/rsca/internal/register"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Client is a api.RSCAClient for co-ordinating requests from the server.
type Client struct {
	Logger   *slog.Logger
	hostname string
	checks   checks.Checks
	inbox    chan *api.Message
	outbox   chan *api.Message
	register *register.Message
}

// NewClient returns a setup api.RSCAClient.
func NewClient(logger *slog.Logger, hostName string, checkList checks.Checks) *Client {
	return &Client{
		Logger:   logger,
		hostname: hostName,
		checks:   checkList,
		inbox:    make(chan *api.Message),
		outbox:   make(chan *api.Message),
	}
}

func (c *Client) streamMessages(ctx context.Context, cancel context.CancelFunc, stream api.RSCA_PipeClient) {
	for {
		in, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			c.Logger.DebugContext(ctx, "EOF found, closing channel")
			cancel()

			return
		} else if s, ok := status.FromError(err); err != nil && ok {
			if s.Code() == codes.Unavailable {
				c.Logger.WarnContext(ctx, "server has gone away", slogtool.ErrorAttr(err))
				cancel()

				return
			}
		} else if err != nil {
			c.Logger.ErrorContext(ctx, "failed to receive a note", slog.Any("msg", in), slogtool.ErrorAttr(err))
			cancel()

			return
		}

		if ctx.Err() != nil {
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

	go c.streamMessages(ctx, cancel, stream)

	return func() error {
		for {
			select {
			case <-ctx.Done():
				c.Logger.DebugContext(ctx, "context cancelled")
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
						c.Logger.ErrorContext(ctx, "unable to send message", slogtool.ErrorAttr(err))
						cancel()
					}
				}
			}
		}
	}
}

// processUpdateAll processes a trigger all message.
func (c *Client) processUpdateAll(ctx context.Context) {
	c.Logger.DebugContext(ctx, "processUpdateAll() called")
	c.checks.NextRun(time.Time{})
}

// processRepeatRegister processes a repeat-registration request message.
func (c *Client) processRepeatRegister(ctx context.Context) {
	c.Logger.DebugContext(ctx, "processRepeatRegister() called")
	c.SendRepeatRegistration(ctx)
}

// SendRepeatRegistration sends the registration message to the server.
func (c *Client) SendRepeatRegistration(ctx context.Context) {
	c.Logger.DebugContext(ctx, "sending repeat registration message")
	c.register.UpdateInfoStat(ctx)

	c.outbox <- api.Message_builder{
		Envelope:            api.Envelope_builder{Sender: c.register.Member(), Recipient: api.MembersByID("_server")}.Build(),
		MemberUpdateMessage: c.register.UpdateMessage(),
	}.Build()
}

// // Send adds a message to the outbox to be sent, may block if channel is full.
// func (c *Client) Send(msg *api.Message) {
// 	c.inbox <- msg
// }

func (c *Client) wrapEventMessage(in *api.EventMessage) *api.Message {
	return api.Message_builder{
		Envelope: api.Envelope_builder{
			Recipient: api.MembersByID("_server"),
			Sender:    c.register.Member(),
		}.Build(),
		EventMessage: in,
	}.Build()
}

// RunEvents runs as a go routine that processes the response channel and creates messages to add to the outbox.
//
//nolint:gocognit // I don't see an easy way to make this less complex without making it less maintainable.
func (c *Client) RunEvents(
	ctx context.Context,
	cancel context.CancelFunc,
	regmsg *register.Message,
	respChan chan *api.EventMessage,
) func() error {
	c.register = regmsg

	return func() error {
		for {
			select {
			case <-ctx.Done():
				c.Logger.DebugContext(ctx, "context cancelled")

				return nil
			case in, ok := <-respChan:
				if ok {
					c.outbox <- c.wrapEventMessage(in)
				}
			case in, ok := <-c.inbox:
				if ok {
					if in != nil {
						switch v := in.WhichMessage(); v { //nolint:exhaustive // default catches unhandled.
						case api.Message_PingMessage_case:
							go func() {
								c.outbox <- common.GeneratePingMessage(ctx, c.Logger, c.hostname, in, in.GetPingMessage())
							}()
						case api.Message_TriggerAllMessage_case:
							go c.processUpdateAll(ctx)
						case api.Message_RepeatRegistrationMessage_case:
							go c.processRepeatRegister(ctx)
						default:
							c.Logger.InfoContext(ctx,
								"Received unhandled message",
								slog.String("message-type", v.String()),
							)
						}
						c.Logger.DebugContext(ctx, "message processing finished")
					} else {
						c.Logger.DebugContext(ctx, "nil message received")
						cancel()
					}
				}
			}
		}
	}
}
