package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/common"
	"github.com/na4ma4/rsca/internal/state"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ssmChannelSize is the size of the buffered channel for serverStreamMessage.
const ssmChannelSize = 2

// Server is a api.RSCAServer for co-ordinating streams from clients.
type Server struct {
	logger   *zap.Logger
	hostname string
	state    state.State
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
	Stream       api.RSCA_PipeServer
	TriggerClose context.CancelFunc
	Record       *api.Member
}

type serverStreamMessage struct {
	M *api.Message
	E error
}

// NewServer returns a prepared server object.
func NewServer(logger *zap.Logger, st state.State) *Server {
	return &Server{
		logger:  logger,
		streams: map[string]*serverStream{},
		state:   st,
		metric: &metric{
			ActiveConnections: promauto.NewGauge(prometheus.GaugeOpts{
				Name:      "connections_active",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "Number of active connections",
			}),
			LifetimeConnections: promauto.NewCounter(prometheus.CounterOpts{
				Name:      "connections_lifetime_total",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "Number of connections (lifetime)",
			}),
			// Received: map[string]*prometheus.CounterVec{},
			Received: promauto.NewCounterVec(prometheus.CounterOpts{
				Name:      "events_received_total",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "received packets, grouped by event",
			}, []string{"source", "event"}),
			EventStatus: promauto.NewCounterVec(prometheus.CounterOpts{
				Name:      "check_results_total",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "received check results",
			}, []string{"source", "check", "result"}),
			PingTick: promauto.NewCounter(prometheus.CounterOpts{
				Name:      "ping_tick_total",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "number of server ticks received",
			}),
			PingMessages: promauto.NewCounter(prometheus.CounterOpts{
				Name:      "ping_messages_total",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "number of tick messages sent",
			}),
			PingMessageErrors: promauto.NewCounter(prometheus.CounterOpts{
				Name:      "ping_message_errors_total",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "number of tick messages that failed to send",
			}),
			PingLatency: promauto.NewGaugeVec(prometheus.GaugeOpts{
				Name:      "ping_latency",
				Namespace: "rsca",
				Subsystem: "server",
				Help:      "ping latency in ms",
			}, []string{"source"}),
		},
	}
}

// TriggerInfo triggers an information update from the host (repeat-registration).
func (s *Server) TriggerInfo(ctx context.Context, m *api.Members) (*api.TriggerInfoResponse, error) {
	msg := &api.Message{
		Envelope: &api.Envelope{Sender: &api.Member{Id: "master"}, Recipient: m},
		Message: &api.Message_RepeatRegistrationMessage{
			RepeatRegistrationMessage: &api.RepeatRegistrationMessage{
				Id: uuid.New().String(),
			},
		},
	}
	if err := s.Send(msg); err != nil {
		s.logger.Error("send returned error", zap.Error(err))

		return nil, status.Error(codes.Internal, err.Error()) //nolint:wrapcheck
	}

	return &api.TriggerInfoResponse{
		Names: s.streamIDsToHostnames(s.streamIDsFromRecipient(m)),
	}, nil
}

// TriggerAll triggers all the services on a matching host.
func (s *Server) TriggerAll(ctx context.Context, m *api.Members) (*api.TriggerAllResponse, error) {
	msg := &api.Message{
		Envelope: &api.Envelope{Sender: &api.Member{Id: "master"}, Recipient: m},
		Message:  &api.Message_TriggerAllMessage{TriggerAllMessage: &api.TriggerAllMessage{Id: uuid.New().String()}},
	}
	if err := s.Send(msg); err != nil {
		s.logger.Error("send returned error", zap.Error(err))

		return nil, status.Error(codes.Internal, err.Error()) //nolint:wrapcheck
	}

	return &api.TriggerAllResponse{
		Names: s.streamIDsToHostnames(s.streamIDsFromRecipient(m)),
	}, nil
}

func (s *Server) streamIDsToHostnames(streamIDs []string) (hostNames []string) {
	for _, streamID := range streamIDs {
		hostNames = append(hostNames, s.streamIDToHostname(streamID))
	}

	return
}

// RemoveHost removes a specified list of hosts from the server.
func (s *Server) RemoveHost(ctx context.Context, in *api.RemoveHostRequest) (*api.RemoveHostResponse, error) {
	s.logger.Debug("RemoveHost()", zap.Strings("targets", in.GetNames()))

	o := &api.RemoveHostResponse{
		Names: []string{},
	}

	for _, hostname := range in.GetNames() {
		if v, ok := s.state.GetMemberByHostname(hostname); ok {
			if streamID, ok := s.state.GetStreamIDByMember(v); ok {
				if st, ok := s.streams[streamID]; ok && st.TriggerClose != nil {
					s.logger.Debug("remove host, closing channel", zap.String("target", hostname), zap.String("streamID", streamID))
					st.TriggerClose()
				}
			}

			if err := s.state.Delete(v); err != nil {
				s.logger.Debug("unable to remove host from state storage", zap.String("target", hostname), zap.Error(err))

				return o, status.Error(codes.Internal, fmt.Sprintf("unable to delete host: %s", err)) //nolint:wrapcheck
			}

			s.logger.Debug("host removed from storage", zap.String("target", hostname))

			o.Names = append(o.Names, v.GetName())
		}

		if v, ok := s.state.GetMemberByHostname(hostname); ok {
			s.logger.Debug("host found in storage after removal", zap.String("target", hostname), zap.Reflect("member", v))
		}
	}

	return o, nil
}

// ListHosts returns a list of hosts currently registered with the server.
func (s *Server) ListHosts(req *api.Empty, stream api.Admin_ListHostsServer) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.state.Walk(func(m *api.Member) error {
		if err := stream.Send(m); err != nil {
			return err //nolint:wrapcheck // goes into status.Error
		}

		return nil
	}); err != nil {
		return status.Error(codes.Internal, err.Error()) //nolint:wrapcheck
	}

	return nil
}

// Pipe handles incoming streams and maintains the stream map.
func (s *Server) Pipe(stream api.RSCA_PipeServer) error {
	streamID := uuid.New().String()
	ctx, cancel := context.WithCancel(context.Background())

	s.lock.Lock()
	s.streams[streamID] = &serverStream{
		Stream:       stream,
		TriggerClose: cancel,
	}
	s.lock.Unlock()

	s.metric.ActiveConnections.Inc()
	s.metric.LifetimeConnections.Inc()

	defer func() {
		s.lock.Lock()
		defer s.lock.Unlock()

		s.logger.Debug("defer delete stream", zap.String("stream.id", streamID))
		s.metric.ActiveConnections.Dec()
		delete(s.streams, streamID)

		_ = s.state.DeactivateByStreamID(streamID)
	}()

	msgStream := s.processPipeMessages(streamID, stream)

	return s.processPipe(ctx, streamID, stream, msgStream)
}

func (s *Server) processPipeMessages(streamID string, stream api.RSCA_PipeServer) chan serverStreamMessage {
	o := make(chan serverStreamMessage, ssmChannelSize)

	go func() {
		for {
			in, err := stream.Recv()
			o <- serverStreamMessage{
				M: in,
				E: err,
			}

			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
					s.logger.Debug("pipe process thread closing stream (with error)",
						zap.String("stream.id", streamID),
						zap.Error(err),
					)

					return
				}

				s.logger.Debug("pipe process thread closing stream", zap.String("stream.id", streamID))

				return
			}
		}
	}()

	return o
}

// processPipe is the main message handler.
//nolint:cyclop // don't see a way to make this much more simpler without making it less readable.
func (s *Server) processPipe(
	ctx context.Context,
	streamID string,
	stream api.RSCA_PipeServer,
	msgStream chan serverStreamMessage,
) error {
	for {
		select {
		case m, ok := <-msgStream:
			if ok {
				if err := m.E; err != nil {
					if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
						s.logger.Debug("closing stream", zap.String("stream.id", streamID), zap.Error(err))

						return nil
					}

					return fmt.Errorf("stream closed: %w", err)
				}

				s.updateLastSeen(streamID, time.Now())

				switch msg := m.M.Message.(type) {
				case *api.Message_EventMessage:
					s.processEventMessage(m.M, msg)
				case *api.Message_RegisterMessage:
					s.processRegisterMessage(streamID, m.M, msg)
				case *api.Message_MemberUpdateMessage:
					s.processMemberUpdateMessage(streamID, m.M, msg)
				case *api.Message_PingMessage:
					s.metric.Received.WithLabelValues("_all", "PingMessage").Inc()
					s.metric.Received.WithLabelValues(m.M.Envelope.Sender.GetName(), "PingMessage").Inc()

					if err := common.ProcessPingMessage(s.logger, stream, s.hostname, m.M, msg); err != nil {
						s.logger.Error("unable to send PongMessage in response to PingMessage", zap.Error(err))
					}
				case *api.Message_PongMessage:
					s.processPongMessage(streamID, m.M, msg)
				default:
					s.metric.Received.WithLabelValues("_all", "Unknown").Inc()
					s.metric.Received.WithLabelValues(m.M.Envelope.Sender.GetName(), "Unknown").Inc()
					s.logger.Info("Received unhandled message", zap.Reflect("message", m.M))
				}
			}
		case <-ctx.Done():
			return fmt.Errorf("stream context closed: %w", ctx.Err())
		}
	}
}

func (s *Server) updateLastSeen(streamID string, t time.Time) {
	s.logger.Debug("updateLastSeen()", zap.String("streamID", streamID))
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.streams[streamID]; ok {
		if v.Record != nil {
			v.Record.LastSeen = timestamppb.New(t)
			v.Record.Active = true
			_ = s.state.AddWithStreamID(streamID, v.Record)
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
		zap.String("source.hostname", in.Envelope.Sender.GetName()), zap.String("check.name", msg.EventMessage.GetCheck()),
		zap.String("check.status", msg.EventMessage.Status.String()),
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
	s.updateMember(streamID, msg.RegisterMessage.Member)
}

func (s *Server) processMemberUpdateMessage(
	streamID string,
	in *api.Message,
	msg *api.Message_MemberUpdateMessage,
) {
	s.metric.Received.WithLabelValues("_all", "MemberUpdateMessage").Inc()
	s.metric.Received.WithLabelValues(in.Envelope.Sender.GetName(), "MemberUpdateMessage").Inc()
	s.logger.Debug("client updated",
		zap.String("rsca.client.name", msg.MemberUpdateMessage.Member.GetName()),
		zap.Strings("rsca.client.tags", msg.MemberUpdateMessage.Member.GetTag()),
		zap.Strings("rsca.client.capabilities", msg.MemberUpdateMessage.Member.GetCapability()),
		zap.Strings("rsca.client.services", msg.MemberUpdateMessage.Member.GetService()),
	)
	s.updateMember(streamID, msg.MemberUpdateMessage.Member)
}

func (s *Server) updateMember(streamID string, m *api.Member) {
	s.logger.Debug("updateMember()", zap.String("streamID", streamID), zap.Reflect("member", m))
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.streams[streamID]; ok {
		m.LastSeen = timestamppb.Now()
		m.Active = true
		s.streams[streamID].Record = m

		_ = s.state.AddWithStreamID(streamID, s.streams[streamID].Record)
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
		s.metric.PingLatency.WithLabelValues(s.streamIDToHostname(streamID)).Set(
			float64(td.Milliseconds()),
		)
		s.setPingLatency(streamID, td)
	}
}

func (s *Server) setPingLatency(streamID string, td time.Duration) {
	s.logger.Debug("setPingLatency()", zap.String("streamID", streamID))
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.streams[streamID]; ok {
		if v.Record != nil {
			v.Record.PingLatency = durationpb.New(td)
			v.Record.Active = true
		}

		_ = s.state.AddWithStreamID(streamID, v.Record)
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
//nolint:cyclop // don't see a way to make this much more simpler without making it less readable.
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
			if stream.Record.IsMatch(r) {
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

	if err := multierr.Combine(errs...); err != nil {
		return fmt.Errorf("unable to send messaage to some or all client streams: %w", err)
	}

	return nil
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
						Ts: timestamppb.Now(),
					}},
				}
				if err := s.Send(msg); err != nil {
					s.logger.Error("send returned error", zap.Error(err))
				}
			}
		}
	}
}
