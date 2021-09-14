package state

import "github.com/na4ma4/rsca/api"

// Member stores a member detail record with annoations that are compatible with asdine/storm.
type Member struct {
	ID       string `storm:"id"`
	StreamID string `storm:"index"`
	Member   *api.Member
}

// func apiMemberToMember(m *api.Member) *Member {
// 	return &Member{
// 		ID:     m.GetName(),
// 		Member: m,
// 	}
// }

// type Members map[string]*Member

// func (a Members) Add(m *api.Member) error {
// 	if m != nil {
// 		a[m.GetName()] = apiMemberToMember(m)
// 	}

// 	return nil
// }

// func (a Members) Walk(walkFunc func(*api.Member) error) error {
// 	for _, v := range a {
// 		if err := walkFunc(v.Member); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (a Members) GetByHostname(hostname string) (*api.Member, bool) {
// 	if v, ok := a[hostname]; ok {
// 		return v.Member, true
// 	}

// 	return nil, false
// }
