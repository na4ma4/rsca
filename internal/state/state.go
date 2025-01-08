package state

import "github.com/na4ma4/rsca/api"

// State is an interface for a service that can store and list the active and historic members.
type State interface {
	// AddWithStreamID adds a member to the internal list along with their streamID.
	AddWithStreamID(streamID string, member *api.Member) error
	// Close will close the underlying storage system or return an error.
	Close() error
	// DeactivateByStreamID sets the Active property on a member to false.
	DeactivateByStreamID(streamID string) error
	// DeactivateByHostname sets the Active property on a member to false by host name.
	DeactivateByHostname(hostName string) error
	// Walk will run a supplied function over each of the members in the storage.
	Walk(walkFunc func(*api.Member) error) error
	// GetMemberByHostname returns a member by their hostname.
	GetMemberByHostname(hostName string) (*api.Member, bool)
	// GetStreamIDByMember returns a stream ID by a specified member.
	GetStreamIDByMember(member *api.Member) (string, bool)
	// Delete removes a member and will disconnect them if they're connected.
	Delete(member *api.Member) error

	// Add(*api.Member) error
	// Deactivate(*api.Member) error
}
