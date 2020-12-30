package gui

import (
	"fmt"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomConfigAssistant struct {
	u                   *gtkUI
	account             *account
	roomConfigComponent *mucRoomConfigComponent
	roomID              jid.Bare
	autoJoin            bool
	currentPageIndex    int
	onSuccess           func(currentAccount *account, roomID jid.Bare, autoJoin bool)
	onCancel            func()

	assistant          gtki.Assistant `gtk-widget:"room-config-assistant"`
	infoPageBox        gtki.Box       `gtk-widget:"room-config-info-page"`
	accessPageBox      gtki.Box       `gtk-widget:"room-config-access-page"`
	permissionsPageBox gtki.Box       `gtk-widget:"room-config-permissions-page"`
	occupantsPageBox   gtki.Box       `gtk-widget:"room-config-occupants-page"`
	othersPageBox      gtki.Box       `gtk-widget:"room-config-others-page"`
	summaryPageBox     gtki.Box       `gtk-widget:"room-config-summary-page"`

	infoPage        mucRoomConfigPage
	accessPage      mucRoomConfigPage
	permissionsPage mucRoomConfigPage
	occupantsPage   mucRoomConfigPage
	othersPage      mucRoomConfigPage
	summaryPage     mucRoomConfigPage

	log coylog.Logger
}

func (u *gtkUI) newRoomConfigAssistant(ca *account, roomID jid.Bare, form *muc.RoomConfigForm, autoJoin bool, onSuccess func(*account, jid.Bare, bool), onCancel func()) *roomConfigAssistant {
	rc := &roomConfigAssistant{
		u:        u,
		account:  ca,
		roomID:   roomID,
		autoJoin: autoJoin,
		log: u.log.WithFields(log.Fields{
			"room":  roomID,
			"where": "configureRoomAssistant",
		}),
	}

	rc.onSuccess = func(a *account, rid jid.Bare, aj bool) {
		if onSuccess != nil {
			onSuccess(a, rid, aj)
		}
	}

	rc.onCancel = func() {
		if onCancel != nil {
			onCancel()
		}
	}

	rc.initBuilder()
	rc.initRoomConfigComponent(form)
	rc.initRoomConfigPages()
	rc.initDefaults()

	return rc
}

func (rc *roomConfigAssistant) initBuilder() {
	b := newBuilder("MUCRoomConfigAssistant")
	panicOnDevError(b.bindObjects(rc))

	b.ConnectSignals(map[string]interface{}{
		"on_page_changed": rc.onPageChanged,
		"on_cancel":       rc.onCancelClicked,
		"on_apply":        rc.onApplyClicked,
	})
}

func (rc *roomConfigAssistant) initRoomConfigComponent(form *muc.RoomConfigForm) {
	rc.roomConfigComponent = rc.u.newMUCRoomConfigComponent(rc.account, rc.roomID, form, rc.autoJoin, rc.assistant)
}

func (rc *roomConfigAssistant) initRoomConfigPages() {
	rc.infoPage = rc.roomConfigComponent.getConfigPage(roomConfigInformationPageIndex)
	rc.accessPage = rc.roomConfigComponent.getConfigPage(roomConfigAccessPageIndex)
	rc.permissionsPage = rc.roomConfigComponent.getConfigPage(roomConfigPermissionsPageIndex)
	rc.occupantsPage = rc.roomConfigComponent.getConfigPage(roomConfigOccupantsPageIndex)
	rc.othersPage = rc.roomConfigComponent.getConfigPage(roomConfigOthersPageIndex)
	rc.summaryPage = rc.roomConfigComponent.getConfigPage(roomConfigSummaryPageIndex)

	rc.infoPageBox.Add(rc.infoPage.getContent())
	rc.accessPageBox.Add(rc.accessPage.getContent())
	rc.permissionsPageBox.Add(rc.permissionsPage.getContent())
	rc.occupantsPageBox.Add(rc.occupantsPage.getContent())
	rc.othersPageBox.Add(rc.othersPage.getContent())
	rc.summaryPageBox.Add(rc.summaryPage.getContent())
}

func (rc *roomConfigAssistant) initDefaults() {
	rc.infoPageBox.SetHExpand(true)
	rc.accessPageBox.SetHExpand(true)
	rc.permissionsPageBox.SetHExpand(true)
	rc.occupantsPageBox.SetHExpand(true)
	rc.othersPageBox.SetHExpand(true)
	rc.summaryPageBox.SetHExpand(true)

	rc.assistant.SetTitle(i18n.Localf("Configuration for room [%s]", rc.roomID.String()))
}

func (rc *roomConfigAssistant) onPageChanged(_ gtki.Assistant, p gtki.Widget) {
	previousPage := rc.pageByIndex(rc.currentPageIndex)
	if !previousPage.isValid() {
		rc.assistant.SetCurrentPage(rc.currentPageIndex)
		return
	}

	previousPage.collectData()

	rc.currentPageIndex = rc.assistant.GetCurrentPage()
	currentPage := rc.pageByIndex(rc.currentPageIndex)

	rc.assistant.SetPageComplete(p, true)
	currentPage.refresh()
}

func (rc *roomConfigAssistant) onCancelClicked() {
	rc.assistant.SetSensitive(false)
	rc.roomConfigComponent.cancelConfiguration(rc.onCancelSuccess, rc.onCancelError)
}

func (rc *roomConfigAssistant) onCancelSuccess() {
	rc.onCancel()
	doInUIThread(rc.assistant.Destroy)
}

func (rc *roomConfigAssistant) onCancelError(err error) {
	// TODO show a friendly error message (bassed on "err") to the user
}

func (rc *roomConfigAssistant) onApplyClicked() {
	rc.assistant.SetSensitive(false)
	rc.summaryPage.showLoadingOverlay()
	rc.roomConfigComponent.submitConfigurationForm(
		rc.onApplySuccess,
		rc.onApplyError,
	)
}

func (rc *roomConfigAssistant) onApplySuccess() {
	rc.onSuccess(rc.account, rc.roomID, rc.roomConfigComponent.autoJoin)
	doInUIThread(rc.assistant.Destroy)
}

func (rc *roomConfigAssistant) onApplyError(err error) {
	// TODO show a friendly error message (bassed on "err") to the user inside the assistant
}

func (rc *roomConfigAssistant) pageByIndex(p int) mucRoomConfigPage {
	switch p {
	case roomConfigInformationPageIndex:
		return rc.infoPage
	case roomConfigAccessPageIndex:
		return rc.accessPage
	case roomConfigPermissionsPageIndex:
		return rc.permissionsPage
	case roomConfigOccupantsPageIndex:
		return rc.occupantsPage
	case roomConfigOthersPageIndex:
		return rc.othersPage
	case roomConfigSummaryPageIndex:
		return rc.summaryPage
	default:
		panic(fmt.Sprintf("developer error: unsupported room config assistant page \"%d\"", p))
	}
}

func (rc *roomConfigAssistant) show() {
	rc.assistant.ShowAll()
}
