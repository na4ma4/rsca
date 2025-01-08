package api

// // MembersByName returns a member from a supplied name.
// func MembersByName(name string) *Members {
// 	return &Members{
// 		Name: []string{name},
// 	}
// }

// MembersByID returns a member from a supplied id.
func MembersByID(id string) *Members {
	return Members_builder{
		Id: []string{id},
	}.Build()
}

// RecipientBySender converts a Envelope.Sender to Envelope.Recipient.
func RecipientBySender(in *Member) *Members {
	if in == nil {
		return &Members{}
	}

	switch {
	case in.GetId() != "":
		return Members_builder{Id: []string{in.GetId()}}.Build()
	case in.GetName() != "":
		return Members_builder{Name: []string{in.GetName()}}.Build()
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
