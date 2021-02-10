package session

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

func newMUCRoomOccupant(nickname string, affiliation data.Affiliation, role data.Role, realJid jid.Full) *muc.Occupant {
	return &muc.Occupant{
		Nickname:    nickname,
		Affiliation: affiliation,
		Role:        role,
		RealJid:     realJid,
	}
}

func (m *mucManager) handleOccupantUpdate(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	l := m.log.WithFields(log.Fields{
		"room":     roomID,
		"occupant": op.Nickname,
		"method":   "handleOccupantUpdate",
	})

	room, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.Error("Trying to get a room that is not in the room manager")
		return
	}

	occupantUpdateInfo := m.newOccupantPresenceUpdateData(room, op)

	updated := room.Roster().UpdateOrAddOccupant(op)
	// Added IsSelfOccupantInTheRoom validation to avoid publishing the events of
	// other occupants until receive the selfPresence.
	// This validation is temporally while 'state machine' pattern is implemented.
	if room.IsSelfOccupantInTheRoom() {
		if updated {
			m.handleOccupantAffiliationRoleUpdate(occupantUpdateInfo)
			m.occupantUpdate(roomID, op)
		} else {
			m.occupantJoined(roomID, op)
		}
	}
}

type occupantPresenceUpdateData struct {
	room                *muc.Room
	currentOccupantInfo *muc.OccupantPresenceInfo
	newOccupantInfo     *muc.OccupantPresenceInfo
	actorOccupant       *data.OccupantUpdateActor
}

func (m *mucManager) occupantPresenceCurrentInfo(room *muc.Room, nickname string) (*muc.OccupantPresenceInfo, error) {
	occupant, exists := room.Roster().GetOccupant(nickname)
	if !exists {
		return nil, errors.New("the occupant is not present in the roster")
	}

	op := &muc.OccupantPresenceInfo{
		Nickname: occupant.Nickname,
		RealJid:  occupant.RealJid,
		AffiliationRole: &muc.OccupantAffiliationRole{
			Affiliation: occupant.Affiliation,
			Role:        occupant.Role,
		},
	}

	return op, nil
}

func (m *mucManager) newOccupantPresenceUpdateData(room *muc.Room, newOccupantInfo *muc.OccupantPresenceInfo) *occupantPresenceUpdateData {
	currentOccupantInfo, err := m.occupantPresenceCurrentInfo(room, newOccupantInfo.Nickname)
	if err != nil {
		m.log.WithError(err).Error("An error occurred when getting the occupant update info")
		return nil
	}

	// Getting the actor affiliation and role
	actorOccupant := &data.OccupantUpdateActor{
		Nickname: newOccupantInfo.AffiliationRole.Actor,
	}

	if actor, ok := room.Roster().GetOccupant(actorOccupant.Nickname); ok {
		actorOccupant.Affiliation = actor.Affiliation
		actorOccupant.Role = actor.Role
	}

	occupantUpdateInfo := &occupantPresenceUpdateData{
		room,
		currentOccupantInfo,
		newOccupantInfo,
		actorOccupant,
	}

	return occupantUpdateInfo
}

func (od *occupantPresenceUpdateData) prevAffiliation() data.Affiliation {
	return od.currentOccupantInfo.AffiliationRole.Affiliation
}

func (od *occupantPresenceUpdateData) newAffiliation() data.Affiliation {
	return od.newOccupantInfo.AffiliationRole.Affiliation
}

func (od *occupantPresenceUpdateData) prevRole() data.Role {
	return od.currentOccupantInfo.AffiliationRole.Role
}

func (od *occupantPresenceUpdateData) newRole() data.Role {
	return od.newOccupantInfo.AffiliationRole.Role
}

func (od *occupantPresenceUpdateData) isOwnPresence() bool {
	return od.currentOccupantInfo.Nickname == od.room.SelfOccupantNickname()
}

func (od *occupantPresenceUpdateData) nickname() string {
	return od.newOccupantInfo.Nickname
}

func (od *occupantPresenceUpdateData) reason() string {
	return od.newOccupantInfo.AffiliationRole.Reason
}

func (m *mucManager) handleOccupantAffiliationRoleUpdate(occupantUpdateInfo *occupantPresenceUpdateData) {
	prevAffiliation := occupantUpdateInfo.prevAffiliation()
	prevRole := occupantUpdateInfo.prevRole()

	newAffiliation := occupantUpdateInfo.newAffiliation()
	newRole := occupantUpdateInfo.newRole()

	switch {
	case !prevAffiliation.Equals(newAffiliation):
		m.handleOccupantAffiliationUpdate(occupantUpdateInfo)

	case prevRole.Name() != newRole.Name():
		m.handleOccupantRoleUpdate(occupantUpdateInfo)
	}
}

func (m *mucManager) handleOccupantAffiliationUpdate(occupantUpdateInfo *occupantPresenceUpdateData) {
	affiliationUpate := data.AffiliationUpdate{
		OccupantUpdateAffiliationRole: data.OccupantUpdateAffiliationRole{
			Actor:    occupantUpdateInfo.actorOccupant,
			Nickname: occupantUpdateInfo.nickname(),
			Reason:   occupantUpdateInfo.reason(),
		},
		New:      occupantUpdateInfo.newAffiliation(),
		Previous: occupantUpdateInfo.prevAffiliation(),
	}

	if occupantUpdateInfo.isOwnPresence() {
		m.selfOccupantAffiliationUpdated(occupantUpdateInfo.room.ID, affiliationUpate)
		return
	}

	m.occupantAffiliationUpdated(occupantUpdateInfo.room.ID, affiliationUpate)
}

func (m *mucManager) handleOccupantRoleUpdate(occupantUpdateInfo *occupantPresenceUpdateData) {
	roleUpdate := data.RoleUpdate{
		OccupantUpdateAffiliationRole: data.OccupantUpdateAffiliationRole{
			Actor:    occupantUpdateInfo.actorOccupant,
			Nickname: occupantUpdateInfo.nickname(),
			Reason:   occupantUpdateInfo.reason(),
		},
		New:      occupantUpdateInfo.newRole(),
		Previous: occupantUpdateInfo.prevRole(),
	}

	if occupantUpdateInfo.isOwnPresence() {
		m.selfOccupantRoleUpdated(occupantUpdateInfo.room.ID, roleUpdate)
		return
	}

	m.occupantRoleUpdated(occupantUpdateInfo.room.ID, roleUpdate)
}

func (m *mucManager) handleOccupantLeft(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	l := m.log.WithFields(log.Fields{
		"room":     roomID,
		"occupant": op.Nickname,
		"method":   "handleOccupantLeft",
	})

	r, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.Error("Trying to get a room that is not in the room manager")
		return
	}

	err := r.Roster().RemoveOccupant(op.Nickname)
	if err != nil {
		l.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
		return
	}

	m.occupantLeft(roomID, op)
}

func (m *mucManager) handleOccupantUnavailable(roomID jid.Bare, op *muc.OccupantPresenceInfo, u *xmppData.MUCUser) {
	if u == nil || u.Destroy == nil {
		return
	}

	m.handleRoomDestroyed(roomID, u.Destroy)
}

func (m *mucManager) handleRoomDestroyed(roomID jid.Bare, d *xmppData.MUCRoomDestroy) {
	j, ok := jid.TryParseBare(d.Jid)
	if d.Jid != "" && !ok {
		m.log.WithFields(log.Fields{
			"room":            roomID,
			"alternativeRoom": d.Jid,
			"method":          "handleRoomDestroyed",
		}).Warn("Invalid alternative room ID")
	}

	m.roomDestroyed(roomID, d.Reason, j, d.Password)
}

func (m *mucManager) handleNonMembersRemoved(roomID jid.Bare, op *muc.OccupantPresenceInfo) {
	l := m.log.WithFields(log.Fields{
		"room":     roomID,
		"occupant": op.Nickname,
		"method":   "handleNonMembersRemoved",
	})

	r, ok := m.roomManager.GetRoom(roomID)
	if !ok {
		l.Error("Trying to get a room that is not in the room manager")
		return
	}

	err := r.Roster().RemoveOccupant(op.Nickname)
	if err != nil {
		l.WithError(err).Error("An error occurred trying to remove the occupant from the roster")
	}

	if r.SelfOccupant().Nickname == op.Nickname {
		m.removeSelfOccupant(roomID)
		_ = m.roomManager.LeaveRoom(roomID)
		return
	}
	m.occupantRemoved(roomID, op.Nickname)
}
