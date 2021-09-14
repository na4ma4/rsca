package state

import (
	"sync"

	"github.com/asdine/storm/v3"
	"github.com/na4ma4/rsca/api"
	"go.uber.org/zap"
)

type Disk struct {
	db     *storm.DB
	logger *zap.Logger
	lock   sync.Mutex
}

func NewDiskState(logger *zap.Logger, filename string) (State, error) {
	db, err := storm.Open(filename)
	if err != nil {
		return nil, err
	}

	return &Disk{
		db:     db,
		logger: logger,
		lock:   sync.Mutex{},
	}, nil
}

func (d *Disk) Close() error {
	return d.db.Close()
}

func (d *Disk) Add(in *api.Member) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	mb := Member{
		ID:     in.GetName(),
		Member: in,
	}

	return d.db.Save(&mb)
}

func (d *Disk) AddWithStreamID(streamID string, in *api.Member) error {
	mb := Member{
		ID:       in.GetName(),
		StreamID: streamID,
		Member:   in,
	}

	return d.db.Save(&mb)
}

func (d *Disk) GetByHostname(hostname string) (*api.Member, bool) {
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

func (d *Disk) Walk(walkFunc func(*api.Member) error) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	var ms []Member
	if err := d.db.All(&ms); err != nil {
		return err
	}

	for _, v := range ms {
		if err := walkFunc(v.Member); err != nil {
			return err
		}
	}

	return nil
}

func (d *Disk) WalkMember(walkFunc func(*Member) error) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	var ms []*Member
	if err := d.db.All(&ms); err != nil {
		return err
	}

	for _, v := range ms {
		if err := walkFunc(v); err != nil {
			return err
		}
	}

	return nil
}

func (d *Disk) Delete(in *api.Member) error {
	var m Member
	if err := d.db.One("ID", in.GetName(), &m); err == nil {
		if err := d.db.DeleteStruct(&m); err != nil {
			return err
		}
	}

	return nil

}

func (d *Disk) Deactivate(in *api.Member) error {
	var m Member
	if err := d.db.One("ID", in.GetName(), &m); err == nil {
		m.Member.Active = false

		if err := d.db.Save(&m); err != nil {
			return err
		}
	}

	return nil
}

func (d *Disk) DeactivateByStreamID(streamID string) error {
	d.logger.Debug("DeactivateByStreamID", zap.String("streamID", streamID))
	var m Member
	if err := d.db.One("StreamID", streamID, &m); err == nil {
		m.Member.Active = false
		m.StreamID = ""

		if err := d.db.Save(&m); err != nil {
			return err
		}
	}

	return nil
}
