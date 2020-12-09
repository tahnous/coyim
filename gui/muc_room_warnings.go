package gui

import (
	"sync"

	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewWarning struct {
	text string
	lock sync.Mutex

	bar     gtki.Box   `gtk-widget:"warning-infobar"`
	message gtki.Label `gtk-widget:"message"`
}

func newRoomViewWarning(text string) *roomViewWarning {
	w := &roomViewWarning{
		text: text,
	}

	builder := newBuilder("MUCRoomWarning")
	panicOnDevError(builder.bindObjects(w))

	w.message.SetText(w.text)

	return w
}

type roomViewWarningsInfoBar struct {
	infoBar gtki.InfoBar `gtk-widget:"bar"`
}

func (v *roomView) newRoomViewWarningsInfoBar() *roomViewWarningsInfoBar {
	ib := &roomViewWarningsInfoBar{}

	builder := newBuilder("MUCRoomWarningsInfoBar")
	panicOnDevError(builder.bindObjects(ib))

	builder.ConnectSignals(map[string]interface{}{
		"on_show_warnings": v.showWarnings,
		"on_close":         v.removeWarningsInfobar,
	})

	return ib
}

// getMessageType implements the "message" interface
func (ib *roomViewWarningsInfoBar) getMessageType() gtki.MessageType {
	return gtki.MESSAGE_WARNING
}

// getWidget implements the "widget" interface
func (ib *roomViewWarningsInfoBar) getWidget() gtki.Widget {
	return ib.infoBar
}

type roomViewWarningsOverlay struct {
	warnings []*roomViewWarning
	onClose  func()

	box      gtki.Box      `gtk-widget:"warningsBox"`
	revealer gtki.Revealer `gtk-widget:"revealer"`
}

func (v *roomView) newRoomViewWarningsOverlay() *roomViewWarningsOverlay {
	o := &roomViewWarningsOverlay{
		onClose: v.closeNotificationsOverlay,
	}

	builder := newBuilder("MUCRoomWarningsOverlay")
	panicOnDevError(builder.bindObjects(o))

	builder.ConnectSignals(map[string]interface{}{
		"on_close": o.close,
	})

	mucStyles.setRoomWarningsBoxStyle(o.box)

	v.messagesBox.Add(o.revealer)

	return o
}

func (o *roomViewWarningsOverlay) add(text string) {
	w := newRoomViewWarning(text)
	o.warnings = append(o.warnings, w)

	mucStyles.setRoomWarningsMessageBoxStyle(w.bar)

	o.box.PackStart(w.bar, false, false, 5)

	o.box.ShowAll()
}

func (o *roomViewWarningsOverlay) show() {
	o.revealer.SetRevealChild(true)
}

func (o *roomViewWarningsOverlay) hide() {
	o.revealer.SetRevealChild(false)
}

func (o *roomViewWarningsOverlay) close() {
	o.hide()
	o.onClose()
}

func (o *roomViewWarningsOverlay) clear() {
	warnings := o.warnings
	for _, w := range warnings {
		o.box.Remove(w.bar)
	}
	o.warnings = nil
}
