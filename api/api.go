package api

// // MembersByName returns a member from a supplied name.
// func MembersByName(name string) *Members {
// 	return &Members{
// 		Name: []string{name},
// 	}
// }

// MembersByID returns a member from a supplied id.
func MembersByID(id string) *Members {
	return &Members{
		Id: []string{id},
	}
}

// RecipientBySender converts a Envelope.Sender to Envelope.Recipient.
func RecipientBySender(in *Member) *Members {
	if in == nil {
		return &Members{}
	}

	switch {
	case in.GetId() != "":
		return &Members{Id: []string{in.GetId()}}
	case in.GetName() != "":
		return &Members{Name: []string{in.GetName()}}
	default:
		return &Members{}
	}
}

// ExitCodeToStatus converts a returned exit code to a Status.
func ExitCodeToStatus(exitCode int) Status {
	m := map[int]Status{
		0: Status_OK,
		1: Status_WARNING,
		2: Status_CRITICAL,
		3: Status_UNKNOWN,
	}

	if v, ok := m[exitCode]; ok {
		return v
	}

	return Status_UNKNOWN
}
