package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/na4ma4/config"
	"github.com/na4ma4/go-slogtool"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/helpers"
	"github.com/na4ma4/rsca/internal/state"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/multierr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ssmChannelSize is the size of the buffered channel for serverStreamMessage.
const ssmChannelSize = 2

// Server is a api.RSCAServer for co-ordinating streams from clients.
type Server struct {
	Logger   *slog.Logger
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
func NewServer(logger *slog.Logger, st state.State) *Server {
	return &Server{
		Logger:  logger,
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
	msg := api.Message_builder{
		Envelope: api.Envelope_builder{
			Sender: api.Member_builder{
				Id: proto.String("master"),
			}.Build(),
			Recipient: m,
		}.Build(),
		RepeatRegistrationMessage: api.RepeatRegistrationMessage_builder{
			Id: proto.String(uuid.New().String()),
		}.Build(),
	}.Build()
	if err := s.Send(ctx, msg); err != nil {
		s.Logger.ErrorContext(ctx, "send returned error", slogtool.ErrorAttr(err))

		return nil, status.Error(codes.Internal, err.Error())
	}

	return api.TriggerInfoResponse_builder{
		Names: s.streamIDsToHostnames(s.streamIDsFromRecipient(m)),
	}.Build(), nil
}

// TriggerAll triggers all the services on a matching host.
func (s *Server) TriggerAll(ctx context.Context, m *api.Members) (*api.TriggerAllResponse, error) {
	msg := api.Message_builder{
		Envelope: api.Envelope_builder{
			Sender: api.Member_builder{
				Id: proto.String("master"),
			}.Build(),
			Recipient: m,
		}.Build(),
		TriggerAllMessage: api.TriggerAllMessage_builder{
			Id: proto.String(uuid.New().String()),
		}.Build(),
	}.Build()
	if err := s.Send(ctx, msg); err != nil {
		s.Logger.ErrorContext(ctx, "send returned error", slogtool.ErrorAttr(err))

		return nil, status.Error(codes.Internal, err.Error())
	}

	return api.TriggerAllResponse_builder{
		Names: s.streamIDsToHostnames(s.streamIDsFromRecipient(m)),
	}.Build(), nil
}

func (s *Server) streamIDsToHostnames(streamIDs []string) []string {
	hostNames := []string{}

	for _, streamID := range streamIDs {
		hostNames = append(hostNames, s.streamIDToHostname(streamID))
	}

	return hostNames
}

// RemoveHost removes a specified list of hosts from the server.
func (s *Server) RemoveHost(ctx context.Context, in *api.RemoveHostRequest) (*api.RemoveHostResponse, error) {
	s.Logger.DebugContext(ctx, "RemoveHost()", slog.Any("targets", in.GetNames()))

	out := []string{}
	o := api.RemoveHostResponse_builder{
		Names: out,
	}.Build()

	for _, hostname := range in.GetNames() {
		if v, ok := s.state.GetMemberByHostname(hostname); ok {
			if streamID, streamIDOK := s.state.GetStreamIDByMember(v); streamIDOK {
				if st, streamOK := s.streams[streamID]; streamOK && st.TriggerClose != nil {
					s.Logger.DebugContext(ctx, "remove host, closing channel",
						slog.String("target", hostname), slog.String("streamID", streamID),
					)
					st.TriggerClose()
				}
			}

			if err := s.state.Delete(v); err != nil {
				s.Logger.DebugContext(ctx, "unable to remove host from state storage",
					slog.String("target", hostname), slogtool.ErrorAttr(err),
				)

				return o, status.Error(codes.Internal, fmt.Sprintf("unable to delete host: %s", err))
			}

			s.Logger.DebugContext(ctx, "host removed from storage", slog.String("target", hostname))

			out = append(out, v.GetName())
		}

		if v, ok := s.state.GetMemberByHostname(hostname); ok {
			s.Logger.DebugContext(ctx, "host found in storage after removal",
				slog.String("target", hostname), slog.Any("member", v),
			)
		}
	}

	o.SetNames(out)

	return o, nil
}

// ListHosts returns a list of hosts currently registered with the server.
func (s *Server) ListHosts(_ *api.Empty, stream api.Admin_ListHostsServer) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := s.state.Walk(func(m *api.Member) error {
		return stream.Send(m)
	}); err != nil {
		return status.Error(codes.Internal, err.Error())
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

		s.Logger.DebugContext(ctx, "defer delete stream", slog.String("stream.id", streamID))
		s.metric.ActiveConnections.Dec()
		delete(s.streams, streamID)

		_ = s.state.DeactivateByStreamID(streamID)
	}()

	msgStream := s.processPipeMessages(ctx, streamID, stream)

	return s.processPipe(ctx, streamID, stream, msgStream)
}

func (s *Server) processPipeMessages(
	ctx context.Context,
	streamID string,
	stream api.RSCA_PipeServer,
) chan serverStreamMessage {
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
					s.Logger.DebugContext(ctx, "pipe process thread closing stream (with error)",
						slog.String("stream.id", streamID),
						slogtool.ErrorAttr(err),
					)

					return
				}

				s.Logger.DebugContext(ctx, "pipe process thread closing stream", slog.String("stream.id", streamID))

				return
			}
		}
	}()

	return o
}

// processPipe is the main message handler.
//
//nolint:gocognit // I don't see an easy way to make this less complex without making it less maintainable.
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
						s.Logger.DebugContext(ctx,
							"closing stream",
							slog.String("stream.id", streamID),
							slogtool.ErrorAttr(err),
						)

						return nil
					}

					return fmt.Errorf("stream closed: %w", err)
				}

				s.updateLastSeen(ctx, streamID, time.Now())

				switch v := m.M.WhichMessage(); v { //nolint:exhaustive // default catches unhandled.
				case api.Message_EventMessage_case:
					s.processEventMessage(ctx, m.M, m.M.GetEventMessage())
				case api.Message_RegisterMessage_case:
					s.processRegisterMessage(ctx, streamID, m.M, m.M.GetRegisterMessage())
				case api.Message_MemberUpdateMessage_case:
					s.processMemberUpdateMessage(ctx, streamID, m.M, m.M.GetMemberUpdateMessage())
				case api.Message_PingMessage_case:
					s.metric.Received.WithLabelValues("_all", "PingMessage").Inc()
					s.metric.Received.WithLabelValues(m.M.GetEnvelope().GetSender().GetName(), "PingMessage").Inc()

					if err := helpers.ProcessPingMessage(ctx, s.Logger, stream, s.hostname, m.M, m.M.GetPingMessage()); err != nil {
						s.Logger.ErrorContext(ctx,
							"unable to send PongMessage in response to PingMessage",
							slogtool.ErrorAttr(err),
						)
					}
				case api.Message_PongMessage_case:
					s.processPongMessage(ctx, streamID, m.M, m.M.GetPongMessage())
				default:
					s.metric.Received.WithLabelValues("_all", "Unknown").Inc()
					s.metric.Received.WithLabelValues(m.M.GetEnvelope().GetSender().GetName(), "Unknown").Inc()
					s.Logger.InfoContext(ctx, "Received unhandled message",
						slog.String("message-type", m.M.WhichMessage().String()),
						slog.Any("message", m.M),
					)
				}
			}
		case <-ctx.Done():
			return fmt.Errorf("stream context closed: %w", ctx.Err())
		}
	}
}

func (s *Server) updateLastSeen(
	ctx context.Context,
	streamID string,
	t time.Time,
) {
	s.Logger.DebugContext(ctx, "updateLastSeen()", slog.String("streamID", streamID))
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.streams[streamID]; ok {
		if v.Record != nil {
			v.Record.SetLastSeen(timestamppb.New(t))
			v.Record.SetActive(true)
			_ = s.state.AddWithStreamID(streamID, v.Record)
		}
	}
}

func (s *Server) processEventMessage(
	ctx context.Context,
	in *api.Message,
	msg *api.EventMessage,
) {
	s.metric.Received.WithLabelValues("_all", "EventMessage").Inc()
	s.metric.Received.WithLabelValues(in.GetEnvelope().GetSender().GetName(), "EventMessage").Inc()
	s.metric.EventStatus.WithLabelValues(
		in.GetEnvelope().GetSender().GetName(),
		msg.GetCheck(),
		msg.GetStatus().String(),
	).Inc()
	s.Logger.DebugContext(ctx, "Received EventMessage")
	s.Logger.InfoContext(ctx, "received check data", slog.String("response.id", msg.GetId()),
		slog.String("source.hostname", in.GetEnvelope().GetSender().GetName()),
		slog.String("check.name", msg.GetCheck()),
		slog.String("check.status", msg.GetStatus().String()),
		slog.String("check.output", msg.GetOutput()))

	if err := writeCheckResponse(ctx, s.Logger, msg); err != nil {
		s.Logger.ErrorContext(ctx, "unable to write check response", slogtool.ErrorAttr(err))
	}
}

func (s *Server) processRegisterMessage(
	ctx context.Context,
	streamID string,
	in *api.Message,
	msg *api.RegisterMessage,
) {
	s.metric.Received.WithLabelValues("_all", "RegisterMessage").Inc()
	s.metric.Received.WithLabelValues(in.GetEnvelope().GetSender().GetName(), "RegisterMessage").Inc()
	s.Logger.InfoContext(ctx, "client registered",
		slog.String("rsca.client.name", msg.GetMember().GetName()),
		slog.Any("rsca.client.tags", msg.GetMember().GetTag()),
		slog.Any("rsca.client.capabilities", msg.GetMember().GetCapability()),
		slog.Any("rsca.client.services", msg.GetMember().GetService()),
	)
	s.updateMember(ctx, streamID, msg.GetMember())
}

func (s *Server) processMemberUpdateMessage(
	ctx context.Context,
	streamID string,
	in *api.Message,
	msg *api.MemberUpdateMessage,
) {
	s.metric.Received.WithLabelValues("_all", "MemberUpdateMessage").Inc()
	s.metric.Received.WithLabelValues(in.GetEnvelope().GetSender().GetName(), "MemberUpdateMessage").Inc()
	s.Logger.DebugContext(ctx, "client updated",
		slog.String("rsca.client.name", msg.GetMember().GetName()),
		slog.Any("rsca.client.tags", msg.GetMember().GetTag()),
		slog.Any("rsca.client.capabilities", msg.GetMember().GetCapability()),
		slog.Any("rsca.client.services", msg.GetMember().GetService()),
	)
	s.updateMember(ctx, streamID, msg.GetMember())
}

func (s *Server) updateMember(ctx context.Context, streamID string, m *api.Member) {
	s.Logger.DebugContext(ctx, "updateMember()", slog.String("streamID", streamID), slog.Any("member", m))
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.streams[streamID]; ok {
		m.SetLastSeen(timestamppb.Now())
		m.SetActive(true)
		s.streams[streamID].Record = m

		_ = s.state.AddWithStreamID(streamID, s.streams[streamID].Record)
	}
}

func (s *Server) processPongMessage(
	ctx context.Context,
	streamID string,
	in *api.Message,
	msg *api.PongMessage,
) {
	s.metric.Received.WithLabelValues("_all", "PongMessage").Inc()
	s.metric.Received.WithLabelValues(in.GetEnvelope().GetSender().GetName(), "PongMessage").Inc()
	s.Logger.DebugContext(ctx, "Received PongMessage",
		slog.String("streamID", streamID),
		slog.String("ping.id", msg.GetId()),
	)

	if v := msg.GetTs(); v != nil {
		td := time.Since(v.AsTime())
		s.metric.PingLatency.WithLabelValues(s.streamIDToHostname(streamID)).Set(
			float64(td.Milliseconds()),
		)
		s.setPingLatency(ctx, streamID, td)
	}
}

func (s *Server) setPingLatency(
	ctx context.Context,
	streamID string,
	td time.Duration,
) {
	s.Logger.DebugContext(ctx, "setPingLatency()", slog.String("streamID", streamID))
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.streams[streamID]; ok {
		if v.Record != nil {
			v.Record.SetPingLatency(durationpb.New(td))
			v.Record.SetActive(true)
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
//
//nolint:gocognit // I don't see an easy way to make this less complex without making it less maintainable.
func (s *Server) streamIDsFromRecipient(in *api.Members) []string {
	s.lock.Lock()
	defer s.lock.Unlock()

	streamIDs := make(map[string]struct{})

	// s.Logger.DebugContext(ctx, "streamIDsFromRecipient()", slog.Any("registered.streams", s.streams))

	for streamID, stream := range s.streams {
		if stream.Record == nil {
			continue
		}

		for _, r := range in.GetId() {
			if strings.EqualFold(stream.Record.GetId(), r) {
				streamIDs[streamID] = struct{}{}
			}
		}

		for _, r := range in.GetName() {
			if stream.Record.IsMatch(r) {
				streamIDs[streamID] = struct{}{}
			}
		}

		if s.compareSlices(stream.Record.GetCapability(), in.GetCapability()) {
			streamIDs[streamID] = struct{}{}
		}

		if stream.Record.GetTag() != nil || s.compareSlices(append(stream.Record.GetTag(), "_all"), in.GetTag()) {
			streamIDs[streamID] = struct{}{}
		}

		if s.compareSlices(stream.Record.GetService(), in.GetService()) {
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
func (s *Server) Send(
	ctx context.Context,
	msg *api.Message,
) error {
	s.Logger.DebugContext(ctx, "Send", slog.Any("msg", msg))
	streamIDs := s.streamIDsFromRecipient(msg.GetEnvelope().GetRecipient())
	s.Logger.DebugContext(ctx, "Send Streams", slog.Any("streamIDs", streamIDs))

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
				s.Logger.DebugContext(ctx, "server.api:Run() context done", slogtool.ErrorAttr(ctx.Err()))

				return nil
			case t := <-ticker.C:
				s.Logger.DebugContext(ctx, "Tick", slog.Time("tick", t))
				s.metric.PingTick.Inc()

				msg := api.Message_builder{
					Envelope: api.Envelope_builder{
						Sender:    api.Member_builder{Id: proto.String("master")}.Build(),
						Recipient: api.Members_builder{Tag: []string{"_all"}}.Build(),
					}.Build(),
					PingMessage: api.PingMessage_builder{
						Id: proto.String(uuid.New().String()),
						Ts: timestamppb.Now(),
					}.Build(),
				}.Build()
				if err := s.Send(ctx, msg); err != nil {
					s.Logger.ErrorContext(ctx, "send returned error", slogtool.ErrorAttr(err))
				}
			}
		}
	}
}
