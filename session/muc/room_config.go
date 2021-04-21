package muc

import (
	"strconv"
	"strings"

	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

const (
	confiFieldFormType              = "http://jabber.org/protocol/muc#roomconfig"
	configFieldRoomName             = "muc#roomconfig_roomname"
	configFieldRoomDescription      = "muc#roomconfig_roomdesc"
	configFieldEnableLogging        = "muc#roomconfig_enablelogging"
	configFieldEnableArchiving      = "muc#roomconfig_enablearchiving"
	configFieldMemberList           = "muc#roomconfig_getmemberlist"
	configFieldLanguage             = "muc#roomconfig_lang"
	configFieldPubsub               = "muc#roomconfig_pubsub"
	configFieldCanChangeSubject     = "muc#roomconfig_changesubject"
	configFieldAllowInvites         = "muc#roomconfig_allowinvites"
	configFieldAllowMemberInvites   = "{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites"
	configFieldAllowPM              = "muc#roomconfig_allowpm"
	configFieldAllowPrivateMessages = "allow_private_messages"
	configFieldMaxOccupantsNumber   = "muc#roomconfig_maxusers"
	configFieldIsPublic             = "muc#roomconfig_publicroom"
	configFieldIsPersistent         = "muc#roomconfig_persistentroom"
	configFieldPresenceBroadcast    = "muc#roomconfig_presencebroadcast"
	configFieldModerated            = "muc#roomconfig_moderatedroom"
	configFieldMembersOnly          = "muc#roomconfig_membersonly"
	configFieldPasswordProtected    = "muc#roomconfig_passwordprotectedroom"
	configFieldPassword             = "muc#roomconfig_roomsecret"
	configFieldOwners               = "muc#roomconfig_roomowners"
	configFieldWhoIs                = "muc#roomconfig_whois"
	configFieldMaxHistoryFetch      = "muc#maxhistoryfetch"
	configFieldMaxHistoryLength     = "muc#roomconfig_historylength"
	configFieldRoomAdmins           = "muc#roomconfig_roomadmins"
)

// RoomConfigForm represents a room configuration form
type RoomConfigForm struct {
	MaxHistoryFetch                ConfigListSingleField
	AllowPrivateMessages           ConfigListSingleField
	OccupantsCanInvite             bool
	OccupantsCanChangeSubject      bool
	Logged                         bool
	RetrieveMembersList            ConfigListMultiField
	Language                       string
	AssociatedPublishSubscribeNode string
	MaxOccupantsNumber             ConfigListSingleField
	MembersOnly                    bool
	Moderated                      bool
	PasswordProtected              bool
	Persistent                     bool
	PresenceBroadcast              ConfigListMultiField
	Public                         bool
	Admins                         []jid.Any
	Description                    string
	Title                          string
	Owners                         []jid.Any
	Password                       string
	Whois                          ConfigListSingleField

	UnknowFields []*RoomConfigFormField
}

// NewRoomConfigForm creates a new room configuration form instance
func NewRoomConfigForm(form *xmppData.Form) *RoomConfigForm {
	cf := &RoomConfigForm{}

	cf.MaxHistoryFetch = newConfigListSingleField([]string{
		RoomConfigOption10,
		RoomConfigOption20,
		RoomConfigOption30,
		RoomConfigOption50,
		RoomConfigOption100,
		RoomConfigOptionNone,
	})

	cf.AllowPrivateMessages = newConfigListSingleField([]string{
		RoomConfigOptionParticipant,
		RoomConfigOptionModerators,
		RoomConfigOptionNone,
	})

	cf.RetrieveMembersList = newConfigListMultiField([]string{
		RoomConfigOptionModerator,
		RoomConfigOptionParticipant,
		RoomConfigOptionVisitor,
	})

	cf.MaxOccupantsNumber = newConfigListSingleField([]string{
		RoomConfigOption10,
		RoomConfigOption20,
		RoomConfigOption30,
		RoomConfigOption50,
		RoomConfigOption100,
		RoomConfigOptionNone,
	})

	cf.PresenceBroadcast = newConfigListMultiField([]string{
		RoomConfigOptionModerator,
		RoomConfigOptionParticipant,
		RoomConfigOptionVisitor,
	})

	cf.Whois = newConfigListSingleField([]string{
		RoomConfigOptionModerators,
		RoomConfigOptionAnyone,
	})

	cf.SetFormFields(form)

	return cf
}

// GetFormData returns a representation of the room config FORM_TYPE as described in the
// XMPP specification for MUC
//
// For more information see:
// https://xmpp.org/extensions/xep-0045.html#createroom-reserved
// https://xmpp.org/extensions/xep-0045.html#example-163
func (rcf *RoomConfigForm) GetFormData() *xmppData.Form {
	fields := map[string][]string{
		"FORM_TYPE":                     {confiFieldFormType},
		configFieldRoomName:             {rcf.Title},
		configFieldRoomDescription:      {rcf.Description},
		configFieldEnableLogging:        {strconv.FormatBool(rcf.Logged)},
		configFieldEnableArchiving:      {strconv.FormatBool(rcf.Logged)},
		configFieldCanChangeSubject:     {strconv.FormatBool(rcf.OccupantsCanChangeSubject)},
		configFieldAllowInvites:         {strconv.FormatBool(rcf.OccupantsCanInvite)},
		configFieldAllowMemberInvites:   {strconv.FormatBool(rcf.OccupantsCanInvite)},
		configFieldAllowPM:              {rcf.AllowPrivateMessages.CurrentValue()},
		configFieldAllowPrivateMessages: {rcf.AllowPrivateMessages.CurrentValue()},
		configFieldMaxOccupantsNumber:   {rcf.MaxOccupantsNumber.CurrentValue()},
		configFieldIsPublic:             {strconv.FormatBool(rcf.Public)},
		configFieldIsPersistent:         {strconv.FormatBool(rcf.Persistent)},
		configFieldModerated:            {strconv.FormatBool(rcf.Moderated)},
		configFieldMembersOnly:          {strconv.FormatBool(rcf.MembersOnly)},
		configFieldPasswordProtected:    {strconv.FormatBool(rcf.PasswordProtected)},
		configFieldPassword:             {rcf.Password},
		configFieldWhoIs:                {rcf.Whois.CurrentValue()},
		configFieldMaxHistoryFetch:      {rcf.MaxHistoryFetch.CurrentValue()},
		configFieldMaxHistoryLength:     {rcf.MaxHistoryFetch.CurrentValue()},
		configFieldLanguage:             {rcf.Language},
		configFieldRoomAdmins:           jidListToStringList(rcf.Admins),
	}

	formFields := []xmppData.FormFieldX{}
	for name, values := range fields {
		formFields = append(formFields, xmppData.FormFieldX{
			Var:    name,
			Values: values,
		})
	}

	return &xmppData.Form{
		Type:   "submit",
		Fields: formFields,
	}
}

// SetFormFields extract the form fields and updates the room config form properties based on each data
func (rcf *RoomConfigForm) SetFormFields(form *xmppData.Form) {
	for _, field := range form.Fields {
		rcf.setField(field)
	}
}

func (rcf *RoomConfigForm) setField(field xmppData.FormFieldX) {
	switch field.Var {
	case configFieldMaxHistoryFetch, configFieldMaxHistoryLength:
		rcf.MaxHistoryFetch.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case configFieldAllowPM, configFieldAllowPrivateMessages:
		rcf.AllowPrivateMessages.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case configFieldAllowInvites, configFieldAllowMemberInvites:
		rcf.OccupantsCanInvite = formFieldBool(field.Values)

	case configFieldCanChangeSubject:
		rcf.OccupantsCanChangeSubject = formFieldBool(field.Values)

	case configFieldEnableLogging, configFieldEnableArchiving:
		rcf.Logged = formFieldBool(field.Values)

	case configFieldMemberList:
		rcf.RetrieveMembersList.UpdateField(field.Values, formFieldOptionsValues(field.Options))

	case configFieldLanguage:
		rcf.Language = formFieldSingleString(field.Values)

	case configFieldPubsub:
		rcf.AssociatedPublishSubscribeNode = formFieldSingleString(field.Values)

	case configFieldMaxOccupantsNumber:
		rcf.MaxOccupantsNumber.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	case configFieldMembersOnly:
		rcf.MembersOnly = formFieldBool(field.Values)

	case configFieldModerated:
		rcf.Moderated = formFieldBool(field.Values)

	case configFieldPasswordProtected:
		rcf.PasswordProtected = formFieldBool(field.Values)

	case configFieldIsPersistent:
		rcf.Persistent = formFieldBool(field.Values)

	case configFieldPresenceBroadcast:
		rcf.PresenceBroadcast.UpdateField(field.Values, formFieldOptionsValues(field.Options))

	case configFieldIsPublic:
		rcf.Public = formFieldBool(field.Values)

	case configFieldRoomAdmins:
		rcf.Admins = formFieldJidList(field.Values)

	case configFieldRoomDescription:
		rcf.Description = formFieldSingleString(field.Values)

	case configFieldRoomName:
		rcf.Title = formFieldSingleString(field.Values)

	case configFieldOwners:
		rcf.Owners = formFieldJidList(field.Values)

	case configFieldPassword:
		rcf.Password = formFieldSingleString(field.Values)

	case configFieldWhoIs:
		rcf.Whois.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))

	default:
		rcf.setUnknowField(field)
	}
}

func (rcf *RoomConfigForm) setUnknowField(field xmppData.FormFieldX) {
	rcf.UnknowFields = append(rcf.UnknowFields, roomConfigFormFieldFactory(field))
}

func roomConfigFormFieldFactory(field xmppData.FormFieldX) *RoomConfigFormField {
	f := &RoomConfigFormField{
		Name:        field.Var,
		Type:        field.Type,
		Label:       field.Label,
		Description: field.Desc,
	}

	switch field.Type {
	case RoomConfigFieldText, RoomConfigFieldTextPrivate:
		f.Value = formFieldSingleString(field.Values)

	case RoomConfigFieldTextMulti:
		f.Value = strings.Join(field.Values, "\n")

	case RoomConfigFieldBoolean:
		f.Value = formFieldBool(field.Values)

	case RoomConfigFieldList:
		ls := newConfigListSingleField(nil)
		ls.UpdateField(formFieldSingleString(field.Values), formFieldOptionsValues(field.Options))
		f.Value = ls

	case RoomConfigFieldListMulti:
		lm := newConfigListMultiField(nil)
		lm.UpdateField(field.Values, formFieldOptionsValues(field.Options))
		f.Value = lm

	case RoomConfigFieldJidMulti:
		f.Value = formFieldJidList(field.Values)

	default:
		f.Value = field.Values
	}

	return f
}

func formFieldBool(values []string) bool {
	return len(values) > 0 && (strings.ToLower(values[0]) == "true" || values[0] == "1")
}

func formFieldSingleString(values []string) string {
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func formFieldOptionsValues(options []xmppData.FormFieldOptionX) (list []string) {
	for _, o := range options {
		list = append(list, o.Value)
	}

	return list
}

func formFieldJidList(values []string) (list []jid.Any) {
	for _, v := range values {
		list = append(list, jid.Parse(v))
	}
	return list
}

func jidListToStringList(jidList []jid.Any) (result []string) {
	for _, j := range jidList {
		result = append(result, j.String())
	}
	return result
}
