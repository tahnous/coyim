package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

func (v *roomView) onDestroyRoom() {
	d := v.newRoomDestroyView(v.window)
	d.show()
}

type roomDestroyView struct {
	room *roomView

	transient            gtki.Window
	dialog               gtki.Dialog `gtk-widget:"destroy-room-dialog"`
	reasonEntry          gtki.Entry  `gtk-widget:"destroy-room-reason-entry"`
	alternativeRoomEntry gtki.Entry  `gtk-widget:"destroy-room-alternative-room-entry"`
	destroyRoomButton    gtki.Button `gtk-widget:"destroy-room-button"`
	spinnerBox           gtki.Box    `gtk-widget:"destroy-room-spinner-box"`
	notificationBox      gtki.Box    `gtk-widget:"notification-area"`

	spinner      *spinner
	notification *notifications

	cancelChannel chan bool
}

func (v *roomView) newRoomDestroyView(t gtki.Window) *roomDestroyView {
	d := &roomDestroyView{
		room:      v,
		transient: t,
	}

	d.initBuilder()
	d.initDefaults(v.u)

	return d
}

func (d *roomDestroyView) initBuilder() {
	builder := newBuilder("MUCRoomDestroyDialog")
	panicOnDevError(builder.bindObjects(d))

	builder.ConnectSignals(map[string]interface{}{
		"on_destroy_clicked":  d.onDestroyRoom,
		"on_cancel_clicked":   d.onCancel,
		"on_dialog_destroyed": d.onDialogDestroy,
	})
}

func (d *roomDestroyView) initDefaults(u *gtkUI) {
	d.dialog.SetTransientFor(d.transient)

	d.spinner = newSpinner()
	d.spinnerBox.Add(d.spinner.getWidget())

	d.notification = u.newNotifications(d.notificationBox)
}

// onDestroyRoom MUST be called from the UI thread
func (d *roomDestroyView) onDestroyRoom() {
	d.disableFieldsAndShowSpinner()

	reason := d.getReason()

	alternativeID, valid := d.getAlternativeRoomID()
	if !valid {
		d.notification.error(i18n.Local("You must type a valid alternative room address for destroying the room."))
		d.enableFieldsAndHideSpinner()
		return
	}

	d.room.tryDestroyRoom(alternativeID, reason, d.onDestroySuccess, d.onDestroyFails)
}

// onDestroySuccess MUST NOT be called from the UI thread
func (d *roomDestroyView) onDestroySuccess() {
	doInUIThread(d.close)
}

// onDestroyFails MUST NOT be called from the UI thread
func (d *roomDestroyView) onDestroyFails(err error) {
	doInUIThread(func() {
		d.enableFields()
		d.notification.error(d.getFriendlyErrorMessage(err))
	})
}

func (d *roomDestroyView) getFriendlyErrorMessage(err error) string {
	switch err {
	case session.ErrDestroyRoomInvalidIQResponse:
		return i18n.Local("We were able to connect to the room service, " +
			"but we received an invalid response from it. Please try again.")
	case session.ErrDestroyRoomForbidden:
		return i18n.Local("You don't have the permission to destroy this room. " +
			"Please get in contact with one of the room owners.")
	case session.ErrDestroyRoomUnknown:
		return i18n.Local("The room's service responded with an unknow error, " +
			"so, the room can't be destroyed. Please try again.")
	case session.ErrDestroyRoomNoResult:
		return i18n.Local("We were able to send the request to destroy the room, " +
			"but the service responded with an unknow result. Please contact the " +
			"room's administrator.")
	default:
		return i18n.Local("An error occurred while destroying the room, please try again.")
	}
}

// onCancel MUST be called from the UI thread
func (d *roomDestroyView) onCancel() {
	d.cancelActiveRequest()

	d.spinner.hide()
	d.close()
}

// onDialogDestroy MUST be called from the UI thread
func (d *roomDestroyView) onDialogDestroy() {
	d.cancelActiveRequest()
}

// cancelActiveRequest MUST be called from the UI thread
func (d *roomDestroyView) cancelActiveRequest() {
	if d.cancelChannel != nil {
		d.cancelChannel <- true
	}
}

// getReason MUST be called from the UI thread
func (d *roomDestroyView) getReason() string {
	t, _ := d.reasonEntry.GetText()
	return t
}

// getAlternativeRoomID MUST be called from the UI thread
//
// This should be "alternative venue" as the protocol says, but
// we prefer to use "alternative room id" in this context
// in order to have a better understanding of what this field means
func (d *roomDestroyView) getAlternativeRoomID() (jid.Bare, bool) {
	t, _ := d.alternativeRoomEntry.GetText()
	if t != "" {
		return jid.TryParseBare(t)
	}
	return nil, true
}

// disableFields MUST be called from the UI thread
func (d *roomDestroyView) disableFields() {
	d.setSensitivityForAllFields(false)
}

// enableFields MUST be called from the UI thread
func (d *roomDestroyView) enableFields() {
	d.setSensitivityForAllFields(true)
}

// setSensitivityForAllFields MUST be called from the UI thread
func (d *roomDestroyView) setSensitivityForAllFields(v bool) {
	d.reasonEntry.SetSensitive(v)
	d.alternativeRoomEntry.SetSensitive(v)
	d.destroyRoomButton.SetSensitive(v)
}

// show MUST be called from the UI thread
func (d *roomDestroyView) show() {
	d.dialog.Show()
}

// close MUST be called from the UI thread
func (d *roomDestroyView) close() {
	d.dialog.Destroy()
}

// disableFieldsAndShowSpinner MUST be called from the UI thread
func (d *roomDestroyView) disableFieldsAndShowSpinner() {
	d.disableFields()
	d.spinner.show()
}

// enableFieldsAndHideSpinner MUST be called from the UI thread
func (d *roomDestroyView) enableFieldsAndHideSpinner() {
	d.enableFields()
	d.spinner.hide()
}
