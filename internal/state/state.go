package state

import "github.com/na4ma4/rsca/api"

type State interface {
	// Add(*api.Member) error
	AddWithStreamID(string, *api.Member) error
	Close() error
	// Delete(*api.Member) error
	// Deactivate(*api.Member) error
	DeactivateByStreamID(streamID string) error
	GetByHostname(string) (*api.Member, bool)
	Walk(walkFunc func(*api.Member) error) error
}
