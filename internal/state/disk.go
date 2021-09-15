package state

import (
	"fmt"
	"sync"

	"github.com/asdine/storm/v3"
	"github.com/na4ma4/rsca/api"
	"go.uber.org/zap"
)

// Disk is a disk based storage that is compatible with the State interface.
type Disk struct {
	db     *storm.DB
	logger *zap.Logger
	lock   sync.Mutex
}

// NewDiskState returns a Disk service that is compatible with the State interface.
func NewDiskState(logger *zap.Logger, filename string) (State, error) {
	db, err := storm.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	return &Disk{
		db:     db,
		logger: logger,
		lock:   sync.Mutex{},
	}, nil
}

// Close will close the underlying storage system or return an error.
func (d *Disk) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("unable to close database: %w", err)
	}

	return nil
}

// func (d *Disk) Add(in *api.Member) error {
// 	d.lock.Lock()
// 	defer d.lock.Unlock()

// 	mb := Member{
// 		ID:     in.GetName(),
// 		Member: in,
// 	}

// 	if err := d.db.Save(&mb); err != nil {
// 		return fmt.Errorf("unable to add member: %w", err)
// 	}

// 	return nil
// }

// AddWithStreamID adds a member to the internal list along with their streamID.
func (d *Disk) AddWithStreamID(streamID string, in *api.Member) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	mb := Member{
		ID:       in.GetName(),
		StreamID: streamID,
		Member:   in,
	}

	if err := d.db.Save(&mb); err != nil {
		return fmt.Errorf("unable to add with streamID: %w", err)
	}

	return nil
}

// GetMemberByHostname returns a member by their hostname.
func (d *Disk) GetMemberByHostname(hostname string) (*api.Member, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()

	var m Member
	if err := d.db.One("ID", hostname, &m); err != nil {
		return nil, false
	}

	if m.Member != nil {
		return m.Member, true
	}

	return nil, false
}

// GetStreamIDByMember returns a stream ID by a specified member.
func (d *Disk) GetStreamIDByMember(in *api.Member) (string, bool) {
	d.lock.Lock()
	defer d.lock.Unlock()

	var m Member
	if err := d.db.One("ID", in.GetName(), &m); err != nil {
		return "", false
	}

	if m.StreamID != "" {
		return m.StreamID, true
	}

	return "", false
}

// Walk will run a supplied function over each of the members in the storage.
func (d *Disk) Walk(walkFunc func(*api.Member) error) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	var ms []Member

	if err := d.db.All(&ms); err != nil {
		return fmt.Errorf("unable to retrieve members: %w", err)
	}

	for _, v := range ms {
		if err := walkFunc(v.Member); err != nil {
			return err
		}
	}

	return nil
}

// func (d *Disk) WalkMember(walkFunc func(*Member) error) error {
// 	d.lock.Lock()
// 	defer d.lock.Unlock()

// 	var ms []*Member

// 	if err := d.db.All(&ms); err != nil {
// 		return fmt.Errorf("unable to retrieve members: %w", err)
// 	}

// 	for _, v := range ms {
// 		if err := walkFunc(v); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// Delete removes a member and will disconnect them if they're connected.
func (d *Disk) Delete(in *api.Member) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	var m Member

	if err := d.db.One("ID", in.GetName(), &m); err == nil {
		if err := d.db.DeleteStruct(&m); err != nil {
			return fmt.Errorf("unable to delete member: %w", err)
		}
	}

	return nil
}

// func (d *Disk) Deactivate(in *api.Member) error {
// 	var m Member

// 	if err := d.db.One("ID", in.GetName(), &m); err == nil {
// 		m.Member.Active = false

// 		if err := d.db.Save(&m); err != nil {
// 			return fmt.Errorf("unable to deactivate by member: %w", err)
// 		}
// 	}

// 	return nil
// }

// DeactivateByStreamID sets the Active property on a member to false.
func (d *Disk) DeactivateByStreamID(streamID string) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	var m Member

	if err := d.db.One("StreamID", streamID, &m); err == nil {
		m.Member.Active = false
		m.StreamID = ""

		if err := d.db.Save(&m); err != nil {
			return fmt.Errorf("unable to deactivate by streamID: %w", err)
		}
	}

	return nil
}
