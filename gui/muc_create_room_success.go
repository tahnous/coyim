package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoomSuccess struct {
	ac         *account
	roomID     jid.Bare
	onJoinRoom func(*account, jid.Bare)

	view gtki.Box `gtk-widget:"createRoomSuccess"`
}

func (v *createMUCRoom) newCreateRoomSuccess() *createMUCRoomSuccess {
	s := &createMUCRoomSuccess{
		onJoinRoom: v.joinRoom,
	}

	s.initBuilder(v)

	// TODO: Same as with the form, we are modifying the parent from the child
	// constructor, which feels a bit confusing
	v.showSuccessView = func(a *account, roomID jid.Bare) {
		v.form.reset()
		v.container.Remove(v.form.view)
		v.success.updateInfo(a, roomID)
		v.container.Add(v.success.view)
	}

	return s
}

func (s *createMUCRoomSuccess) initBuilder(v *createMUCRoom) {
	builder := newBuilder("MUCCreateRoomSuccess")
	panicOnDevError(builder.bindObjects(s))

	builder.ConnectSignals(map[string]interface{}{
		"on_createRoom_clicked": v.showCreateForm,
		"on_joinRoom_clicked": func() {
			s.onJoinRoom(s.ac, s.roomID)
		},
	})
}

func (s *createMUCRoomSuccess) updateInfo(a *account, roomID jid.Bare) {
	s.ac = a
	s.roomID = roomID
}

func (s *createMUCRoomSuccess) reset() {
	s.ac = nil
	s.roomID = nil
}
