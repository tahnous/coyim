package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/gotk3adapter/gtki"
)

// This is a component that can be used from other views in order to make it easy
// to "retry" any asynchronous operation that might have a "success" and a "failure".
// In this case, this component cover the "failure" action.

type dialogErrorComponent struct {
	title   string
	header  string
	message string

	dialog       gtki.Dialog `gtk-widget:"room-error-dialog"`
	errorTitle   gtki.Label  `gtk-widget:"room-error-dialog-title"`
	errorMessage gtki.Label  `gtk-widget:"room-error-dialog-message"`

	// retry is a callback that will be asynchronously executed when the user wants to
	// "retry" the failed operation for in which this component was used
	retry func()
}

func createDialogErrorComponent(title, header, message string, retry func()) *dialogErrorComponent {
	d := &dialogErrorComponent{
		title:   title,
		header:  header,
		message: message,
		retry:   retry,
	}

	d.initBuilder()
	d.initDefaults()

	return d
}

func (d *dialogErrorComponent) initDefaults() {
	mucStyles.setLabelBoldStyle(d.errorTitle)

	d.dialog.SetTitle(d.title)
	d.errorTitle.SetText(d.header)
	d.errorMessage.SetText(d.message)
}

func (d *dialogErrorComponent) initBuilder() {
	builder := newBuilder("MUCRoomDialogErrorComponent")
	panicOnDevError(builder.bindObjects(d))

	builder.ConnectSignals(map[string]interface{}{
		"on_cancel": d.onCancel,
		"on_retry":  d.onRetry,
	})
}

func (d *dialogErrorComponent) onCancel() {
	d.dialog.Destroy()
}

func (d *dialogErrorComponent) onRetry() {
	go d.retry()
	d.dialog.Destroy()
}

func (d *dialogErrorComponent) show() {
	d.dialog.Show()
}

func (d *dialogErrorComponent) updateMessageForDestroyError(err error) {
	msg := ""
	switch err {
	case session.ErrDestroyRoomInvalidIQResponse, session.ErrDestroyRoomNoResult:
		msg = i18n.Local("We were able to connect to the room service, " +
			"but we received an invalid response from it. Please try again later.")
	case session.ErrDestroyRoomForbidden:
		msg = i18n.Local("You don't have the permission to destroy this room. " +
			"Please contact one of the room owners.")
	case session.ErrDestroyRoomDoesntExist:
		msg = i18n.Local("We couldn't find the room.")
	default:
		msg = i18n.Local("An unknown error occurred during the process. Please try again later.")
	}

	d.errorMessage.SetText(msg)
}
