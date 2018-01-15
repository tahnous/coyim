package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/otr_client"
	"github.com/coyim/coyim/session/access"
)

type executable interface {
	execute(u *gtkUI)
}

type connectAccountCmd struct{ a *account }
type disconnectAccountCmd struct{ a *account }
type connectionInfoCmd struct{ a *account }
type editAccountCmd struct{ a *account }
type removeAccountCmd struct{ a *account }
type toggleAutoConnectCmd struct{ a *account }
type toggleAlwaysEncryptCmd struct{ a *account }

func (u *gtkUI) ExecuteCmd(c interface{}) {
	u.commands <- c
}

func (c connectAccountCmd) execute(u *gtkUI) {
	doInUIThread(func() {
		u.connectAccount(c.a)
	})
}

func (c disconnectAccountCmd) execute(u *gtkUI) {
	go c.a.session.Close()
}

func (c connectionInfoCmd) execute(u *gtkUI) {
	doInUIThread(func() {
		u.connectionInfoDialog(c.a)
	})
}

func (c editAccountCmd) execute(u *gtkUI) {
	doInUIThread(func() {
		u.editAccount(c.a)
	})
}

func (c removeAccountCmd) execute(u *gtkUI) {
	doInUIThread(func() {
		u.removeAccount(c.a)
	})
}

func (c toggleAutoConnectCmd) execute(u *gtkUI) {
	go u.toggleAutoConnectAccount(c.a)
}

func (c toggleAlwaysEncryptCmd) execute(u *gtkUI) {
	go u.toggleAlwaysEncryptAccount(c.a)
}

func (u *gtkUI) watchCommands() {
	for command := range u.commands {
		switch c := command.(type) {
		case executable:
			c.execute(u)
		case otr_client.AuthorizeFingerprintCmd:
			account := c.Account
			uid := c.Peer
			fpr := c.Fingerprint

			//TODO: it could be a different pointer,
			//find the account by ID()
			account.AuthorizeFingerprint(uid.Representation(), fpr)
			u.ExecuteCmd(otr_client.SaveApplicationConfigCmd{})

			ac := u.findAccountForSession(c.Session.(access.Session))
			if ac != nil {
				peer := c.Peer
				convWindowNowOrLater(ac, peer, u, func(cv conversationView) {
					cv.displayNotification(i18n.Localf("You have verified the identity of %s.", peer))
				})
			}
		case otr_client.SaveInstanceTagCmd:
			account := c.Account
			account.InstanceTag = c.InstanceTag
			u.ExecuteCmd(otr_client.SaveApplicationConfigCmd{})
		case otr_client.SaveApplicationConfigCmd:
			u.SaveConfig()
		}
	}
}
