package client

import (
	"context"
	"io"
	"time"

	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/checks"
	"github.com/na4ma4/rsca/internal/common"
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

		if err == io.EOF {
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

// Pipe processes a stream and the inbox received from the server.
func (c *Client) Pipe(ctx context.Context, cancel context.CancelFunc, stream api.RSCA_PipeClient) func() error {
	go c.streamMessages(cancel, stream)

	return func() error {
		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("context cancelled")

				return nil
			case out := <-c.outbox:
				if err := stream.Send(out); err != nil {
					c.logger.Error("unable to send PongMessage in response to PingMessage", zap.Error(err))
				}
			case in := <-c.inbox:
				switch msg := in.Message.(type) {
				case *api.Message_PingMessage:
					if err := common.ProcessPingMessage(c.logger, stream, c.hostname, in, msg); err != nil {
						c.logger.Error("unable to send PongMessage in response to PingMessage", zap.Error(err))
					}
				case *api.Message_UpdateAllMessage:
					c.processUpdateAll()
				default:
					c.logger.Info("Received unhandled message", zap.Reflect("message", in))
				}
				c.logger.Debug("message processing finished")
			}
		}
	}
}

// Pipe processes a stream and the inbox received from the server.
func (c *Client) processUpdateAll() {
	c.checks.NextRun(time.Time{})
}

// // Send adds a message to the outbox to be sent, may block if channel is full.
// func (c *Client) Send(msg *api.Message) {
// 	c.inbox <- msg
// }

// RunEvents runs as a go routine that processes the response channel and creates messages to add to the outbox.
func (c *Client) RunEvents(ctx context.Context, ms *api.Member, respChan chan *api.EventMessage) func() error {
	return func() error {
		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("context cancelled")

				return nil
			case in := <-respChan:
				c.outbox <- &api.Message{
					Envelope: &api.Envelope{
						Recipient: api.MembersByID("_server"),
						Sender:    ms,
					},
					Message: &api.Message_EventMessage{
						EventMessage: in,
					},
				}
			}
		}
	}
}
