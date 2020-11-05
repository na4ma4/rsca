package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is a api.RSCAServer for co-ordinating streams from clients.
type Server struct {
	logger   *zap.Logger
	hostname string
	streams  map[string]*serverStream
	lock     sync.Mutex
	metric   *metric
}

type metric struct {
	ActiveConnections   prometheus.Gauge
	LifetimeConnections prometheus.Counter
	Received            *prometheus.CounterVec
	PingTick            prometheus.Counter
	PingMessages        prometheus.Counter
	PingMessageErrors   prometheus.Counter
	EventStatus         *prometheus.CounterVec
	PingLatency         *prometheus.GaugeVec
}

type serverStream struct {
	Stream api.RSCA_PipeServer
	Record *api.Member
}

// NewServer returns a prepared server object.
func NewServer(logger *zap.Logger, hostName string) *Server {
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
			// Received: map[string]*prometheus.CounterVec{},
			Received: promauto.NewCounterVec(prometheus.CounterOpts{
				Name:      "events_received",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "received packets, grouped by event",
			}, []string{"source", "event"}),
			EventStatus: promauto.NewCounterVec(prometheus.CounterOpts{
				Name:      "check_results",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "received check results",
			}, []string{"source", "check", "result"}),
			PingTick: promauto.NewCounter(prometheus.CounterOpts{
				Name:      "ping_tick",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "number of server ticks received",
			}),
			PingMessages: promauto.NewCounter(prometheus.CounterOpts{
				Name:      "ping_messages",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "number of tick messages sent",
			}),
			PingMessageErrors: promauto.NewCounter(prometheus.CounterOpts{
				Name:      "ping_message_errors",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "number of tick messages that failed to send",
			}),
			PingLatency: promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name:      "ping_latency_ms",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "ping latency in ms",
			}, []string{"source"}),
		},
	}
}

// TriggerAll triggers all the services on a matching host.
func (s *Server) TriggerAll(ctx context.Context, m *api.Members) (*api.TriggerAllResponse, error) {
	msg := &api.Message{
		Envelope: &api.Envelope{Sender: &api.Member{Id: "master"}, Recipient: m},
		Message:  &api.Message_UpdateAllMessage{UpdateAllMessage: &api.UpdateAllMessage{Id: uuid.New().String()}},
	}
	if err := s.Send(msg); err != nil {
		s.logger.Error("send returned error", zap.Error(err))

		return nil, status.Error(codes.Internal, err.Error())
	}

	streamIDs := s.streamIDsFromRecipient(m)
	hostNames := []string{}

	for _, streamID := range streamIDs {
		hostNames = append(hostNames, s.streamIDToHostname(streamID))
	}

	return &api.TriggerAllResponse{
		Names: hostNames,
	}, nil
}

// ListHosts returns a list of hosts currently registered with the server.
func (s *Server) ListHosts(req *api.Empty, stream api.Admin_ListHostsServer) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, m := range s.streams {
		if err := stream.Send(m.Record); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	return nil
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

	return s.processPipe(streamID, stream)
}

// processPipe is the main message handler.
func (s *Server) processPipe(streamID string, stream api.RSCA_PipeServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
				s.logger.Debug("closing stream", zap.String("stream.id", streamID), zap.Error(err))

				return nil
			}

			return fmt.Errorf("stream closed: %w", err)
		}

		s.updateLastSeen(streamID, time.Now())

		switch msg := in.Message.(type) {
		case *api.Message_EventMessage:
			s.processEventMessage(in, msg)
		case *api.Message_RegisterMessage:
			s.processRegisterMessage(streamID, in, msg)
		case *api.Message_PingMessage:
			s.metric.Received.WithLabelValues("_all", "PingMessage").Inc()
			s.metric.Received.WithLabelValues(in.Envelope.Sender.GetName(), "PingMessage").Inc()

			if err := common.ProcessPingMessage(s.logger, stream, s.hostname, in, msg); err != nil {
				s.logger.Error("unable to send PongMessage in response to PingMessage", zap.Error(err))
			}
		case *api.Message_PongMessage:
			s.processPongMessage(streamID, in, msg)
		default:
			s.metric.Received.WithLabelValues("_all", "Unknown").Inc()
			s.metric.Received.WithLabelValues(in.Envelope.Sender.GetName(), "Unknown").Inc()
			s.logger.Info("Received unhandled message", zap.Reflect("message", in))
		}
	}
}

func (s *Server) updateLastSeen(streamID string, t time.Time) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.streams[streamID]; ok {
		if v.Record != nil {
			ts, err := ptypes.TimestampProto(t)
			if err != nil {
				return
			}

			v.Record.LastSeen = ts
		}
	}
}

func (s *Server) processEventMessage(
	in *api.Message,
	msg *api.Message_EventMessage,
) {
	s.metric.Received.WithLabelValues("_all", "EventMessage").Inc()
	s.metric.Received.WithLabelValues(in.Envelope.Sender.GetName(), "EventMessage").Inc()
	s.metric.EventStatus.WithLabelValues(
		in.Envelope.Sender.GetName(),
		msg.EventMessage.GetCheck(),
		msg.EventMessage.GetStatus().String(),
	).Inc()
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
	in *api.Message,
	msg *api.Message_RegisterMessage,
) {
	s.metric.Received.WithLabelValues("_all", "RegisterMessage").Inc()
	s.metric.Received.WithLabelValues(in.Envelope.Sender.GetName(), "RegisterMessage").Inc()
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
	in *api.Message,
	msg *api.Message_PongMessage,
) {
	s.metric.Received.WithLabelValues("_all", "PongMessage").Inc()
	s.metric.Received.WithLabelValues(in.Envelope.Sender.GetName(), "PongMessage").Inc()
	s.logger.Debug("Received PongMessage",
		zap.String("streamID", streamID),
		zap.String("ping.id", msg.PongMessage.GetId()),
	)

	if v := msg.PongMessage.GetTs(); v != nil {
		td := time.Since(v.AsTime())
		s.metric.PingLatency.WithLabelValues(s.streamIDToHostname(msg.PongMessage.GetStreamId())).Set(
			float64(td.Milliseconds()),
		)
		s.setPingLatency(streamID, td)
	}
}

func (s *Server) setPingLatency(streamID string, td time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.streams[streamID]; ok {
		if v.Record != nil {
			v.Record.PingLatency = ptypes.DurationProto(td)
		}
	}
}

func (s *Server) streamIDToHostname(streamID string) string {
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.streams[streamID]; ok {
		if v.Record != nil {
			return v.Record.GetName()
		}
	}

	return ""
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

	// s.logger.Debug("streamIDsFromRecipient()", zap.Reflect("registered.streams", s.streams))

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

		if stream.Record.Tag != nil || s.compareSlices(append(stream.Record.Tag, "_all"), in.Tag) {
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

	errs := []error{nil}

	for _, streamID := range streamIDs {
		if v, ok := s.streams[streamID]; ok {
			s.metric.PingMessages.Inc()

			if err := v.Stream.Send(msg); err != nil {
				s.metric.PingMessageErrors.Inc()

				errs = append(errs, err)
			}
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
				s.logger.Debug("server.api:Run() context done", zap.Error(ctx.Err()))

				return nil
			case t := <-ticker.C:
				s.logger.Debug("Tick", zap.Time("tick", t))
				s.metric.PingTick.Inc()

				msg := &api.Message{
					Envelope: &api.Envelope{Sender: &api.Member{Id: "master"}, Recipient: &api.Members{Tag: []string{"_all"}}},
					Message: &api.Message_PingMessage{PingMessage: &api.PingMessage{
						Id: uuid.New().String(),
						Ts: ptypes.TimestampNow(),
					}},
				}
				if err := s.Send(msg); err != nil {
					s.logger.Error("send returned error", zap.Error(err))
				}
			}
		}
	}
}
