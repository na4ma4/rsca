package server

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

// Server is a api.RSCAServer for co-ordinating streams from clients.
type Server struct {
	logger  *zap.Logger
	streams map[string]*serverStream
	lock    sync.Mutex
	metric  *metric
}

type metric struct {
	ActiveConnections   prometheus.Gauge
	LifetimeConnections prometheus.Counter
}

type serverStream struct {
	Stream api.RSCA_PipeServer
	Record *api.Member
}

// NewServer returns a prepared server object.
func NewServer(logger *zap.Logger) api.RSCAServer {
	return &Server{
		logger:  logger,
		streams: map[string]*serverStream{},
		metric: &metric{
			ActiveConnections: promauto.NewGauge(prometheus.GaugeOpts{
				Name:      "connections_active",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "Number of active connections",
			}),
			LifetimeConnections: promauto.NewCounter(prometheus.CounterOpts{
				Name:      "connections_lifetime",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "Number of connections (lifetime)",
			}),
		},
	}
}

// Pipe handles incoming streams and maintains the stream map.
func (s *Server) Pipe(stream api.RSCA_PipeServer) error {
	streamID := uuid.New().String()

	s.lock.Lock()
	s.streams[streamID] = &serverStream{Stream: stream}
	s.lock.Unlock()

	s.metric.ActiveConnections.Inc()
	s.metric.LifetimeConnections.Inc()

	defer func() {
		s.lock.Lock()
		defer s.lock.Unlock()

		s.logger.Debug("defer delete stream", zap.String("stream.id", streamID))
		s.metric.ActiveConnections.Dec()
		delete(s.streams, streamID)
	}()

	for {
		in, err := stream.Recv()

		switch {
		case err == io.EOF:
			s.logger.Debug("closing stream", zap.String("stream.id", streamID))

			return nil
		case errors.Is(err, context.Canceled):
			s.logger.Debug("closing stream (context cancelled)", zap.String("stream.id", streamID), zap.Error(err))

			return nil
		case err != nil:
			s.logger.Debug("closing stream (error)", zap.String("stream.id", streamID), zap.Error(err))

			return err
		}

		switch msg := in.Message.(type) {
		case *api.Message_EventMessage:
			s.processEventMessage(msg)
		case *api.Message_RegisterMessage:
			s.processRegisterMessage(streamID, msg)
		case *api.Message_PingMessage:
			if err := common.ProcessPingMessage(s.logger, stream, in, msg); err != nil {
				s.logger.Error("unable to send PongMessage in response to PingMessage", zap.Error(err))
			}
		case *api.Message_PongMessage:
			s.processPongMessage(streamID, msg)
		default:
			s.logger.Info("Received unhandled message", zap.Reflect("message", in))
		}
	}
}

func (s *Server) processEventMessage(
	msg *api.Message_EventMessage,
) {
	s.logger.Debug("Received EventMessage")

	s.logger.Info("received check data", zap.String("response.id", msg.EventMessage.GetId()),
		zap.String("check.name", msg.EventMessage.GetCheck()), zap.String("check.status", msg.EventMessage.Status.String()),
		zap.String("check.output", msg.EventMessage.GetOutput()))

	if err := writeCheckResponse(s.logger, msg.EventMessage); err != nil {
		s.logger.Error("unable to write check response", zap.Error(err))
	}
}

func (s *Server) processRegisterMessage(
	streamID string,
	msg *api.Message_RegisterMessage,
) {
	s.logger.Info("client registered",
		zap.String("rsca.client.name", msg.RegisterMessage.Member.GetName()),
		zap.Strings("rsca.client.tags", msg.RegisterMessage.Member.GetTag()),
		zap.Strings("rsca.client.capabilities", msg.RegisterMessage.Member.GetCapability()),
		zap.Strings("rsca.client.services", msg.RegisterMessage.Member.GetService()),
	)
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.streams[streamID]; ok {
		s.streams[streamID].Record = msg.RegisterMessage.Member
	}
}

func (s *Server) processPongMessage(
	streamID string,
	msg *api.Message_PongMessage,
) {
	s.logger.Debug("Received PongMessage",
		zap.String("streamID", streamID),
		zap.String("ping.id", msg.PongMessage.GetId()),
	)
}

func (s *Server) compareSlices(s1, s2 []string) bool {
	if s1 != nil && s2 != nil {
		for _, r := range s1 {
			for _, c := range s2 {
				if strings.EqualFold(c, r) {
					return true
				}
			}
		}
	}

	return false
}

// streamIDsFromRecipient processes a recipient (*api.Members) and returns a list of streamIDs.
func (s *Server) streamIDsFromRecipient(in *api.Members) []string {
	s.lock.Lock()
	defer s.lock.Unlock()

	streamIDs := make(map[string]struct{})

	for streamID, stream := range s.streams {
		if stream.Record == nil {
			continue
		}

		for _, r := range in.Id {
			if strings.EqualFold(stream.Record.Id, r) {
				streamIDs[streamID] = struct{}{}
			}
		}

		for _, r := range in.Name {
			if strings.EqualFold(stream.Record.Name, r) {
				streamIDs[streamID] = struct{}{}
			}
		}

		if s.compareSlices(stream.Record.Capability, in.Capability) {
			streamIDs[streamID] = struct{}{}
		}

		if stream.Record.Tag != nil && s.compareSlices(append(stream.Record.Tag, "_all"), in.Tag) {
			streamIDs[streamID] = struct{}{}
		}

		if s.compareSlices(stream.Record.Service, in.Service) {
			streamIDs[streamID] = struct{}{}
		}
	}

	keys := make([]string, 0, len(streamIDs))
	for k := range streamIDs {
		keys = append(keys, k)
	}

	return keys
}

// Send sends a supplied message to the clients specified in the api.Message:Recipients.
func (s *Server) Send(msg *api.Message) error {
	s.logger.Debug("Send", zap.Reflect("msg", msg))
	streamIDs := s.streamIDsFromRecipient(msg.Envelope.Recipient)
	s.logger.Debug("Send Streams", zap.Strings("streamIDs", streamIDs))

	s.lock.Lock()
	defer s.lock.Unlock()

	var errs []error

	for _, streamID := range streamIDs {
		if v, ok := s.streams[streamID]; ok {
			errs = append(errs, v.Stream.Send(msg))
		}
	}

	return multierr.Combine(errs...)
}

// Run is the background runner for the server.
func (s *Server) Run(ctx context.Context, cfg config.Conf) func() error {
	ticker := time.NewTicker(cfg.GetDuration("server.tick"))

	return func() error {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()

				return nil
			case t := <-ticker.C:
				s.logger.Debug("Tick", zap.Time("tick", t))

				msg := &api.Message{
					Envelope: &api.Envelope{Sender: &api.Member{Id: "master"}, Recipient: &api.Members{Tag: []string{"_all"}}},
					Message:  &api.Message_PingMessage{PingMessage: &api.PingMessage{Id: uuid.New().String()}},
				}
				if err := s.Send(msg); err != nil {
					s.logger.Error("send returned error", zap.Error(err))
				}
			}
		}
	}
}
