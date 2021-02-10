package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
)

func getDisplayRoomSubjectForNickname(nickname, subject string) string {
	if nickname == "" {
		return i18n.Localf("Someone has updated the room subject to: \"%s\"", subject)
	}

	return i18n.Localf("%s updated the room subject to \"%s\"", nickname, subject)
}

func getDisplayRoomSubject(subject string) string {
	if subject == "" {
		return i18n.Local("The room does not have a subject")
	}

	return i18n.Localf("The room subject is \"%s\"", subject)
}

func getDisplayForOccupantAffiliationRoleUpdate(affiliationRoleUpdate data.AffiliationRoleUpdate) string {
	d := newAffiliationRoleUpdateDisplayData(affiliationRoleUpdate)
	return displayAffiliationUpdateMessage(d)
}

func getDisplayForOccupantAffiliationUpdate(affiliationUpdate data.AffiliationUpdate) string {
	d := newAffiliationUpdateDisplayData(affiliationUpdate)
	return displayAffiliationUpdateMessage(d)
}

func getDisplayForSelfOccupantAffiliationUpdate(affiliationUpdate data.AffiliationUpdate) string {
	d := newSelfAffiliationUpdateDisplayData(affiliationUpdate)
	return displayAffiliationUpdateMessage(d)
}

func getDisplayForOccupantRoleUpdate(roleUpdate data.RoleUpdate) string {
	d := newRoleUpdateDisplayData(roleUpdate)
	return displayRoleUpdateMessage(d)
}

type affiliationUpdateDisplayData struct {
	nickname            string
	newAffiliation      data.Affiliation
	previousAffiliation data.Affiliation
	actor               string
	actorAffiliation    data.Affiliation
	reason              string
}

type affiliationRoleUpdateDisplayData struct {
	*affiliationUpdateDisplayData
	newRole      data.Role
	previousRole data.Role
}

type roleUpdateDisplayData struct {
	nickname         string
	newRole          data.Role
	previousRole     data.Role
	actor            string
	actorAffiliation data.Affiliation
	reason           string
}

type affiliationUpdateDisplayer interface {
	affiliation() data.Affiliation
	addExtraMessageInfo() string
	updateReason() string
	previousAffiliationIsNone() bool
	displayForAffiliationRemoved() string
	displayForAffiliationOutcast() string
	displayForAffiliationAdded() string
	displayForAffiliationChanged() string
}

func displayAffiliationUpdateMessage(d affiliationUpdateDisplayer) (message string) {
	newAffiliation := d.affiliation()

	switch {
	case newAffiliation.IsNone():
		message = d.displayForAffiliationRemoved()
	case newAffiliation.IsBanned():
		message = d.displayForAffiliationOutcast()
	default:
		if d.previousAffiliationIsNone() {
			message = d.displayForAffiliationAdded()
		} else {
			message = d.displayForAffiliationChanged()
		}
	}

	if append := d.addExtraMessageInfo(); append != "" {
		message = i18n.Localf("%s %s", message, append)
	}

	if reason := d.updateReason(); reason != "" {
		message = i18n.Localf("%s because: %s", message, reason)
	}

	return message
}

type roleUpdateDisplayer interface {
	displayForRoleChanged() string
	updateReason() string
}

func displayRoleUpdateMessage(d roleUpdateDisplayer) (message string) {
	message = d.displayForRoleChanged()

	if d.updateReason() != "" {
		message = i18n.Localf("%s because: %s", message, d.updateReason())
	}

	return message
}

func newAffiliationRoleUpdateDisplayData(affiliationRoleUpdate data.AffiliationRoleUpdate) *affiliationRoleUpdateDisplayData {
	d := &affiliationRoleUpdateDisplayData{
		affiliationUpdateDisplayData: newAffiliationUpdateDisplayData(affiliationRoleUpdate.AffiliationUpdate),
		newRole:                      affiliationRoleUpdate.RoleUpdate.New,
		previousRole:                 affiliationRoleUpdate.RoleUpdate.Previous,
	}

	if affiliationRoleUpdate.Actor != nil {
		d.actor = affiliationRoleUpdate.Actor.Nickname
		d.actorAffiliation = affiliationRoleUpdate.Actor.Affiliation
	}

	return d
}

func newAffiliationUpdateDisplayData(affiliationUpdate data.AffiliationUpdate) *affiliationUpdateDisplayData {
	d := &affiliationUpdateDisplayData{
		nickname:            affiliationUpdate.Nickname,
		newAffiliation:      affiliationUpdate.New,
		previousAffiliation: affiliationUpdate.Previous,
		reason:              affiliationUpdate.Reason,
	}

	if affiliationUpdate.Actor != nil {
		d.actor = affiliationUpdate.Actor.Nickname
		d.actorAffiliation = affiliationUpdate.Actor.Affiliation
	}

	return d
}

func newRoleUpdateDisplayData(roleUpdate data.RoleUpdate) *roleUpdateDisplayData {
	d := &roleUpdateDisplayData{
		nickname:     roleUpdate.Nickname,
		newRole:      roleUpdate.New,
		previousRole: roleUpdate.Previous,
		reason:       roleUpdate.Reason,
	}

	if roleUpdate.Actor != nil {
		d.actor = roleUpdate.Actor.Nickname
		d.actorAffiliation = roleUpdate.Actor.Affiliation
	}

	return d

}

func (d *affiliationUpdateDisplayData) affiliation() data.Affiliation {
	return d.newAffiliation
}

func (d *affiliationUpdateDisplayData) previousAffiliationIsNone() bool {
	return d.previousAffiliation.IsNone()
}

func (d *affiliationUpdateDisplayData) updateReason() string {
	return d.reason
}

func (d *affiliationUpdateDisplayData) displayForAffiliationRemoved() string {
	if d.actor == "" {
		return i18n.Localf("The %s position of %s was removed",
			displayNameForAffiliation(d.previousAffiliation), d.nickname)
	}
	return i18n.Localf("%s removed the %s position from %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForAffiliation(d.previousAffiliation),
		d.nickname,
	)
}

func displayActorWithAffiliation(actor string, affiliation data.Affiliation) string {
	if affiliation != nil {
		return i18n.Localf("The %s %s", displayNameForAffiliation(affiliation), actor)
	}
	return actor
}

func (d *affiliationUpdateDisplayData) displayForAffiliationOutcast() string {
	if d.actor == "" {
		return i18n.Localf("%s was banned from the room", d.nickname)
	}
	return i18n.Localf("%s banned %s from the room",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		d.nickname,
	)
}

func (d *affiliationUpdateDisplayData) displayForAffiliationAdded() string {
	if d.actor == "" {
		return i18n.Localf("%s is now %s", d.nickname,
			displayNameForAffiliationWithPreposition(d.newAffiliation))
	}
	return i18n.Localf("%s changed the position of %s to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		d.nickname,
		displayNameForAffiliation(d.newAffiliation),
	)
}

func (d *affiliationUpdateDisplayData) displayForAffiliationChanged() string {
	if d.actor == "" {
		return i18n.Localf("The position of %s was changed from %s to %s", d.nickname,
			displayNameForAffiliation(d.previousAffiliation),
			displayNameForAffiliation(d.newAffiliation))
	}
	return i18n.Localf("%s changed the position of %s from %s to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		d.nickname,
		displayNameForAffiliation(d.previousAffiliation),
		displayNameForAffiliation(d.newAffiliation),
	)
}

func (d *affiliationUpdateDisplayData) addExtraMessageInfo() string {
	return ""
}

type selfAffiliationUpdateDisplayData struct {
	*affiliationUpdateDisplayData
}

func newSelfAffiliationUpdateDisplayData(affiliationUpdate data.AffiliationUpdate) *selfAffiliationUpdateDisplayData {
	return &selfAffiliationUpdateDisplayData{
		newAffiliationUpdateDisplayData(affiliationUpdate),
	}
}

func (d *selfAffiliationUpdateDisplayData) displayForAffiliationRemoved() string {
	if d.actor == "" {
		return i18n.Localf("Your position of %s was removed", displayNameForAffiliation(d.previousAffiliation))
	}
	return i18n.Localf("%s removed your position as %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForAffiliation(d.previousAffiliation),
	)
}

func (d *selfAffiliationUpdateDisplayData) displayForAffiliationOutcast() string {
	if d.actor == "" {
		return i18n.Local("You were banned from the room")
	}
	return i18n.Localf("%s banned you from the room",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForAffiliation(d.actorAffiliation),
		d.actor,
	)
}

func (d *selfAffiliationUpdateDisplayData) displayForAffiliationAdded() string {
	if d.actor == "" {
		return i18n.Localf("You are now %s", displayNameForAffiliationWithPreposition(d.newAffiliation))
	}
	return i18n.Localf("%s changed your position to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForAffiliation(d.newAffiliation),
	)
}

func (d *selfAffiliationUpdateDisplayData) displayForAffiliationChanged() string {
	if d.actor == "" {
		return i18n.Localf("Your position was changed from %s to %s",
			displayNameForAffiliation(d.previousAffiliation),
			displayNameForAffiliation(d.newAffiliation))
	}
	return i18n.Localf("%s changed your position from %s to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForAffiliation(d.previousAffiliation),
		displayNameForAffiliation(d.newAffiliation))
}

func (d *selfAffiliationUpdateDisplayData) addExtraMessageInfo() string {
	return ""
}

func displaySelfOccupantAffiliationUpdate(affiliationUpdate data.AffiliationUpdate) string {
	d := newSelfAffiliationUpdateDisplayData(affiliationUpdate)
	return displayAffiliationUpdateMessage(d)
}

func displayNameForAffiliation(affiliation data.Affiliation) string {
	switch {
	case affiliation.IsAdmin():
		return i18n.Local("administrator")
	case affiliation.IsOwner():
		return i18n.Local("owner")
	case affiliation.IsBanned():
		return i18n.Local("outcast")
	case affiliation.IsMember():
		return i18n.Local("member")
	default: // Other values get the default treatment
		return ""
	}
}

func displayNameForAffiliationWithPreposition(affiliation data.Affiliation) string {
	switch {
	case affiliation.IsAdmin():
		return i18n.Local("an admininistrator")
	case affiliation.IsOwner():
		return i18n.Local("an owner")
	case affiliation.IsMember():
		return i18n.Local("a member")
	default: // Other values get the default treatment
		return ""
	}
}

func displayNameForRole(role data.Role) string {
	switch {
	case role.IsModerator():
		return i18n.Local("moderator")
	case role.IsParticipant():
		return i18n.Local("participant")
	default:
		return ""
	}
}

func (d *roleUpdateDisplayData) displayForRoleChanged() string {
	if d.actor == "" {
		return i18n.Localf("The role of %s was changed from %s to %s", d.nickname,
			displayNameForRole(d.previousRole),
			displayNameForRole(d.newRole))
	}
	return i18n.Localf("%s changed the role of %s from %s to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		d.nickname,
		displayNameForRole(d.previousRole),
		displayNameForRole(d.newRole),
	)
}

func (d *roleUpdateDisplayData) updateReason() string {
	return d.reason
}
