package register

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/na4ma4/config"
	"github.com/na4ma4/rsca/api"
	"github.com/na4ma4/rsca/internal/checks"
	"github.com/shirou/gopsutil/v3/host"
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
	hostName, version, buildDate, gitHash string,
	checkList checks.Checks,
	startTime time.Time,
) *Message {
	checkNames := []string{}

	for _, check := range checkList {
		if check.Type == api.CheckType_SERVICE {
			checkNames = append(checkNames, check.Name)
		}
	}

	mb := &api.Member{
		Id:           uuid.New().String(),
		Name:         hostName,
		Capability:   []string{"client", fmt.Sprintf("rsca-%s", version)},
		Service:      checkNames,
		Tag:          cfg.GetStringSlice("general.tags"),
		Version:      version,
		BuildDate:    buildDate,
		GitHash:      gitHash,
		ProcessStart: timestamppb.New(startTime),
	}

	if ut, err := host.BootTimeWithContext(context.Background()); err == nil {
		mb.SystemStart = timestamppb.New(time.Unix(int64(ut), 0))
	}

	if is, err := api.InfoWithContext(context.Background(), time.Now()); err == nil {
		mb.InfoStat = is
	}

	return &Message{
		member: mb,
	}
}

// Message returns the actual api.RegisterMessage.
func (msg *Message) Message() *api.RegisterMessage {
	msg.lock.Lock()
	defer msg.lock.Unlock()

	return &api.RegisterMessage{
		Member: msg.member,
	}
}

// UpdateMessage returns the actual api.MemberUpdateMessage.
func (msg *Message) UpdateMessage() *api.MemberUpdateMessage {
	msg.lock.Lock()
	defer msg.lock.Unlock()

	return &api.MemberUpdateMessage{
		Member: msg.member,
	}
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

	if is, err := api.InfoWithContext(context.Background(), time.Now()); err == nil {
		msg.member.InfoStat = is
	}
}
