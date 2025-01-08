package register

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/dosquad/go-cliversion"
	"github.com/google/uuid"
	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/checks"
	"github.com/shirou/gopsutil/v3/host"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Message is a struct that wraps the information required for an api.RegisterMessage to be sent.
type Message struct {
	member *api.Member
	lock   sync.Mutex
}

// New returns a Message pre-populated.
func New(
	cfg config.Conf,
	hostName string,
	versionInfo *cliversion.VersionInfo,
	checkList checks.Checks,
	startTime time.Time,
) *Message {
	checkNames := []string{}

	for _, check := range checkList {
		if check.Type == api.CheckType_SERVICE {
			checkNames = append(checkNames, check.Name)
		}
	}

	mb := api.Member_builder{
		Id:           proto.String(uuid.New().String()),
		Name:         proto.String(hostName),
		Capability:   []string{"client", "rsca-" + versionInfo.GetBld().GetVersion()},
		Service:      checkNames,
		Tag:          cfg.GetStringSlice("general.tags"),
		Version:      proto.String(versionInfo.GetBld().GetVersion()),
		BuildDate:    proto.String(versionInfo.GetBld().GetDate().AsTime().Format(time.RFC3339)),
		GitHash:      proto.String(versionInfo.GetGit().GetCommit()),
		ProcessStart: timestamppb.New(startTime),
	}.Build()

	if ut, err := host.BootTimeWithContext(context.Background()); err == nil {
		if ut < math.MaxInt64 {
			mb.SetSystemStart(timestamppb.New(time.Unix(int64(ut), 0)))
		}
	}

	if is, err := api.InfoWithContext(context.Background(), time.Now()); err == nil {
		mb.SetInfoStat(is)
	}

	return &Message{
		member: mb,
	}
}

// Message returns the actual api.RegisterMessage.
func (msg *Message) Message() *api.RegisterMessage {
	msg.lock.Lock()
	defer msg.lock.Unlock()

	return api.RegisterMessage_builder{
		Member: msg.member,
	}.Build()
}

// UpdateMessage returns the actual api.MemberUpdateMessage.
func (msg *Message) UpdateMessage() *api.MemberUpdateMessage {
	msg.lock.Lock()
	defer msg.lock.Unlock()

	return api.MemberUpdateMessage_builder{
		Member: msg.member,
	}.Build()
}

// Member returns the member details used in the api.RegisterMessage.
func (msg *Message) Member() *api.Member {
	msg.lock.Lock()
	defer msg.lock.Unlock()

	return msg.member
}

// UpdateInfoStat updates the infostat on the member.
func (msg *Message) UpdateInfoStat(ctx context.Context) {
	msg.lock.Lock()
	defer msg.lock.Unlock()

	if is, err := api.InfoWithContext(ctx, time.Now()); err == nil {
		msg.member.SetInfoStat(is)
	}
}
