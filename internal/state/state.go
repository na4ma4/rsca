package state

import "github.com/na4ma4/rsca/api"

// State is an interface for a service that can store and list the active and historic members.
type State interface {
	// AddWithStreamID adds a member to the internal list along with their streamID.
	AddWithStreamID(string, *api.Member) error
	// Close will close the underlying storage system or return an error.
	Close() error
	// DeactivateByStreamID sets the Active property on a member to false.
	DeactivateByStreamID(streamID string) error
	// Walk will run a supplied function over each of the members in the storage.
	Walk(walkFunc func(*api.Member) error) error

	// GetByHostname returns a member by their hostname.
	// GetByHostname(string) (*api.Member, bool)

	// Add(*api.Member) error
	// Delete(*api.Member) error
	// Deactivate(*api.Member) error
}
