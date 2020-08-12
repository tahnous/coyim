package session

import (
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

const (
	// MUCStatusPresenceJoined inform user that presence refers to one of its own room occupants
	MUCStatusPresenceJoined = "110"
)

func (s *session) isMUCPresence(stanza *data.ClientPresence) bool {
	return stanza.MUC != nil
}

func (s *session) handleMUCPresence(stanza *data.ClientPresence) {
	from := jid.Parse(stanza.From)
	rid, nickname := from.PotentialSplit()

	switch {
	case stanza.MUCUser != nil:
		if stanza.MUCUser.Item != nil {
			s.mucOccupantUpdate(rid.String(), string(nickname), stanza.MUCUser.Item.Affiliation, stanza.MUCUser.Item.Role)
		}

		if len(stanza.MUCUser.Status) > 0 {
			affiliation := stanza.MUCUser.Item.Affiliation
			jid := stanza.MUCUser.Item.Jid
			role := stanza.MUCUser.Item.Role
			for _, status := range stanza.MUCUser.Status {
				switch status.Code {
				case MUCStatusPresenceJoined:
					s.mucOccupantJoined(rid.String(), string(nickname), affiliation, jid, role, status.Code, true)
				}
			}
		}
	}
}

func (s *session) mucOccupantUpdate(rid, nickname, affiliation, role string) {
	s.publishEvent(events.MUCOccupantUpdated{
		MUCOccupant: &events.MUCOccupant{
			MUC: &events.MUC{
				From: rid,
			},
			Nickname: nickname,
		},
		Affiliation: affiliation,
		Role:        role,
	})

	s.mucRosterUpdated()
}

func (s *session) mucRosterUpdated() {
	s.publishEvent(events.MUCOccupantUpdate)
}

func (s *session) mucOccupantJoined(rid, nickname, affiliation, jid, role, status string, v bool) {
	s.publishEvent(events.MUCOccupantJoined{
		MUCOccupantUpdated: &events.MUCOccupantUpdated{
			MUCOccupant: &events.MUCOccupant{
				MUC: &events.MUC{
					From: rid,
				},
				Nickname: nickname,
			},
			Affiliation: affiliation,
			Jid:         jid,
			Role:        role,
			Status:      status,
		},
		Joined: v,
	})
}

func (s *session) hasSomeConferenceService(identities []data.DiscoveryIdentity) bool {
	for _, identity := range identities {
		return identity.Category == "conference" && identity.Type == "text"
	}
	return false
}

func (s *session) filterOnlyChatServices(items *data.DiscoveryItemsQuery) []string {
	chatServices := make([]string, 0)
	for _, item := range items.DiscoveryItems {
		iq, err := s.conn.QueryServiceInformation(item.Jid)
		if err != nil {
			s.log.WithError(err).Error("Error getting the information query for the service:", item.Jid)
			continue
		}
		if iq != nil && s.hasSomeConferenceService(iq.Identities) {
			chatServices = append(chatServices, item.Jid)
		}
	}
	return chatServices
}

//GetChatServices offers the chat services from a xmpp server.
func (s *session) GetChatServices(server jid.Domain) ([]string, error) {
	items, err := s.conn.QueryServiceItems(server.String())
	if err != nil {
		return nil, err
	}
	return s.filterOnlyChatServices(items), nil
}
