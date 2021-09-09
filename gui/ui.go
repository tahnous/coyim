package gui

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/config"
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/gui/settings"
	"github.com/coyim/coyim/i18n"
	ournet "github.com/coyim/coyim/net"
	rosters "github.com/coyim/coyim/roster"
	sessions "github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3adapter/pangoi"
)

const (
	programName        = "CoyIM"
	applicationID      = "im.coy.CoyIM"
	localizationDomain = "coy"
)

type gtkUI struct {
	roster *roster
	app    gtki.Application
	window gtki.ApplicationWindow `gtk-widget:"mainWindow"`

	contactsMenuItem  gtki.MenuItem `gtk-widget:"ContactsMenu"`
	accountsMenuItem  gtki.MenuItem `gtk-widget:"AccountsMenu"`
	chatRoomsMenuItem gtki.MenuItem `gtk-widget:"ChatRoomsMenu"`
	viewMenuItem      gtki.MenuItem `gtk-widget:"ViewMenu"`
	optionsMenuItem   gtki.MenuItem `gtk-widget:"OptionsMenu"`

	searchBox        gtki.Box       `gtk-widget:"search-box"`
	search           gtki.SearchBar `gtk-widget:"search-area"`
	searchEntry      gtki.Entry     `gtk-widget:"search-entry"`
	notificationArea gtki.Box       `gtk-widget:"notification-area"`
	viewMenu         *viewMenu
	optionsMenu      *optionsMenu

	unified       *unifiedLayout
	unifiedCached *unifiedLayout

	internalConfig     *config.ApplicationConfig
	internalConfigLock sync.Mutex

	*accountManager

	displaySettings  *displaySettings
	keyboardSettings *keyboardSettings

	keySupplier config.KeySupplier

	tags *tags

	toggleConnectAllAutomaticallyRequest chan bool
	setShowAdvancedSettingsRequest       chan bool

	commands chan interface{}

	sessionFactory sessions.Factory

	dialerFactory interfaces.DialerFactory

	settings *settings.Settings

	//Desktop notifications
	deNotify    *desktopNotifications
	actionTimes map[string]time.Time

	log coylog.Logger

	hooks OSHooks

	mainBuilder *builder

	ouit *outsideUIThread

	haveConfigEntries *callbacksSet
}

// Graphics represent the graphic configuration
type Graphics struct {
	gtk   gtki.Gtk
	glib  glibi.Glib
	gdk   gdki.Gdk
	pango pangoi.Pango
	extra interface{}
}

// CreateGraphics creates a Graphic represention from the given arguments
func CreateGraphics(gtkVal gtki.Gtk, glibVal glibi.Glib, gdkVal gdki.Gdk, pangoVal pangoi.Pango, extra interface{}) Graphics {
	return Graphics{
		gtk:   gtkVal,
		glib:  glibVal,
		gdk:   gdkVal,
		pango: pangoVal,
		extra: extra,
	}
}

var g Graphics

var coyimVersion string

// UI is the user interface functionality exposed to main
type UI interface {
	Loop()
}

func argsWithApplicationName() *[]string {
	newSlice := make([]string, len(os.Args))
	copy(newSlice, os.Args)
	newSlice[0] = programName
	return &newSlice
}

// NewGTK returns a new client for a GTK ui
func NewGTK(version string, sf sessions.Factory, df interfaces.DialerFactory, gx Graphics, hooks OSHooks, translationDirectory string) UI {
	runtime.LockOSThread()

	inuit := &inUIThread{g: gx}
	outuit := &outsideUIThread{
		doInUIThread: func(f func(*inUIThread)) {
			_ = inuit.g.glib.IdleAdd(func() { f(inuit) })
		},
	}

	coyimVersion = version
	g = gx
	initSignals(inuit)

	registerFinalizerReaping(gx.glib)

	//*.mo files should be in ./i18n/locale_code.utf8/LC_MESSAGES/
	g.glib.InitI18n(localizationDomain, translationDirectory)
	g.gtk.Init(argsWithApplicationName())
	ensureInstalled()

	ret := &gtkUI{
		commands:                             make(chan interface{}, 5),
		toggleConnectAllAutomaticallyRequest: make(chan bool, 100),
		setShowAdvancedSettingsRequest:       make(chan bool, 100),
		dialerFactory:                        df,

		actionTimes: make(map[string]time.Time),
		deNotify:    newDesktopNotifications(),
		log:         log.StandardLogger().WithField("component", "gui"),
		hooks:       hooks,

		ouit: outuit,

		haveConfigEntries: newCallbacksSet(),
	}

	ret.initMUC()

	hooks.AfterInit()

	var err error
	flags := glibi.APPLICATION_FLAGS_NONE
	if *config.MultiFlag {
		flags = glibi.APPLICATION_NON_UNIQUE
	}
	ret.app, err = g.gtk.ApplicationNew(applicationID, flags)
	if err != nil {
		panic(err)
	}

	ret.keySupplier = config.CachingKeySupplier(ret.getMasterPassword)

	ret.accountManager = newAccountManager(ret, ret.log)

	ret.sessionFactory = sf

	ret.settings = settings.For("")

	ret.addAction(ret.app, "quit", ret.quit)
	ret.addAction(ret.app, "about", ret.aboutDialog)
	ret.addAction(ret.app, "preferences", ret.showGlobalPreferences)

	return ret
}

func (u *gtkUI) confirmAccountRemoval(acc *config.Account, removeAccountFunc func(*config.Account)) {
	builder := newBuilder("ConfirmAccountRemoval")

	obj := builder.getObj("RemoveAccount")
	dialog := obj.(gtki.MessageDialog)
	dialog.SetTransientFor(u.window)
	_ = dialog.SetProperty("secondary-text", acc.Account)

	response := dialog.Run()
	if gtki.ResponseType(response) == gtki.RESPONSE_YES {
		removeAccountFunc(acc)
	}

	dialog.Destroy()
}

type torRunningNotification struct {
	area  gtki.Box   `gtk-widget:"infobar"`
	image gtki.Image `gtk-widget:"image"`
	label gtki.Label `gtk-widget:"message"`
}

// TODO: add a spinner
func torRunningNotificationInit(info gtki.Box) *torRunningNotification {
	b := newBuilder("TorRunningNotification")
	torRunningNotif := &torRunningNotification{}
	panicOnDevError(b.bindObjects(torRunningNotif))

	info.Add(torRunningNotif.area)
	torRunningNotif.area.ShowAll()

	return torRunningNotif
}

func (n *torRunningNotification) renderTorNotification(label, imgName string) {
	doInUIThread(func() {
		prov := providerWithCSS("box { background-color: #f1f1f1; color: #000000; border: 1px solid #d3d3d3; border-radius: 2px;}")
		updateWithStyle(n.area, prov)
	})

	n.label.SetText(i18n.Local(label))
	n.image.SetFromIconName(imgName, gtki.ICON_SIZE_BUTTON)
}

func (u *gtkUI) installTor() {
	builder := newBuilder("TorInstallHelper")

	obj := builder.getObj("dialog")
	dialog := obj.(gtki.MessageDialog)
	info := builder.getObj("tor-running-notification").(gtki.Box)
	torNotif := torRunningNotificationInit(info)

	builder.ConnectSignals(map[string]interface{}{
		"on_close": func() {
			dialog.Destroy()
		},
		// TODO: change logos
		"on_press_label": func() {
			if !ournet.Tor.Detect() {
				err := "Tor is still not running"
				torNotif.renderTorNotification(err, "software-update-urgent")
				u.log.Info("Tor is still not running")
			} else {
				err := "Tor is now running"
				torNotif.renderTorNotification(err, "emblem-default")
				u.log.Info("Tor is now running")
			}
		},
	})

	doInUIThread(func() {
		dialog.SetTransientFor(u.window)
		dialog.ShowAll()
	})
}

func (u *gtkUI) wouldYouLikeToInstallTor(k func(bool)) {
	builder := newBuilder("TorHelper")

	dialog := builder.getObj("TorHelper")
	torHelper := dialog.(gtki.MessageDialog)
	torHelper.SetDefaultResponse(gtki.RESPONSE_YES)
	torHelper.SetTransientFor(u.window)

	responseType := gtki.ResponseType(torHelper.Run())
	result := responseType == gtki.RESPONSE_YES
	torHelper.Destroy()
	k(result)
}

func (u *gtkUI) initialSetupWindow() {
	if !ournet.Tor.Detect() {
		u.log.Info("Tor is not running")
		u.wouldYouLikeToInstallTor(func(res bool) {
			if res {
				u.installTor()
			} else {
				u.initialSetupForConfigFile()
			}
		})
	} else {
		u.initialSetupForConfigFile()
	}
}

func (u *gtkUI) initialSetupForConfigFile() {
	u.wouldYouLikeToEncryptYourFile(func(res bool) {
		u.config().SetShouldSaveFileEncrypted(res)
		k := func() {
			go u.showFirstAccountWindow()
		}
		if res {
			u.captureInitialMasterPassword(k, func() {})
		} else {
			k()
		}
	})
}

func (u *gtkUI) loadConfig(configFile string) {
	u.config().WhenLoaded(u.configLoaded)

	ok := false
	var conf *config.ApplicationConfig
	var err error
	for !ok {
		conf, ok, err = config.LoadOrCreate(configFile, u.keySupplier)
		if !ok {
			u.log.WithError(err).Warn("couldn't open encrypted file - either the user didn't supply a password, or the password was incorrect")
			u.keySupplier.Invalidate()
			u.keySupplier.LastAttemptFailed()
		}
	}

	// We assign config here, AFTER the return - so that a nil config means we are in a state of incorrectness and shouldn't do stuff.
	// We never check, since a panic here is a serious programming error
	u.setConfig(conf)

	if err != nil {
		u.log.WithError(err).Warn("something went wrong")
		doInUIThread(u.initialSetupWindow)
		return
	}

	if u.config().UpdateToLatestVersion() {
		_ = u.saveConfigOnlyInternal()
	}
}

func (u *gtkUI) updateUnifiedOrNot() {
	if u.settings.GetSingleWindow() && u.unified == nil {
		u.unified = u.unifiedCached
	}
	if !u.settings.GetSingleWindow() {
		u.unified = nil
	}
}

func (u *gtkUI) configLoaded(c *config.ApplicationConfig) {
	u.settings = settings.For(c.GetUniqueID())
	u.roster.restoreCollapseStatus()
	u.deNotify.updateWith(u.settings)
	u.updateUnifiedOrNot()

	u.buildAccounts(c, u.sessionFactory, u.dialerFactory)

	doInUIThread(func() {
		if u.viewMenu != nil {
			u.viewMenu.setFromConfig(c)
		}

		if u.optionsMenu != nil {
			u.optionsMenu.setFromConfig(c)
		}

		if u.window != nil {
			_, _ = u.window.Emit(accountChangedSignal.String())
		}
	})

	u.addInitialAccountsToRoster()

	if c.ConnectAutomatically {
		u.connectAllAutomatics(false)
	}

	go u.listenToToggleConnectAllAutomatically()
	go u.listenToSetShowAdvancedSettings()
}

func (u *gtkUI) saveConfigInternal() error {
	err := u.saveConfigOnlyInternal()
	if err != nil {
		return err
	}

	u.addNewAccountsFromConfig(u.config(), u.sessionFactory, u.dialerFactory)

	if u.window != nil {
		_, _ = u.window.Emit(accountChangedSignal.String())
	}

	return nil
}

func (u *gtkUI) saveConfigOnlyInternal() error {
	return u.config().Save(u.keySupplier)
}

func (u *gtkUI) SaveConfig() {
	go func() {
		err := u.saveConfigInternal()
		if err != nil {
			u.log.WithError(err).Warn("Failed to save config file")
		}
	}()
}

func (u *gtkUI) removeSaveReload(acc *config.Account) {
	//TODO: the account configs should be managed by the account manager
	u.accountManager.removeAccount(acc, func() {
		u.config().Remove(acc)
		u.SaveConfig()
	})
}

func (u *gtkUI) saveConfigOnly() {
	go func() {
		err := u.saveConfigOnlyInternal()
		if err != nil {
			u.log.WithError(err).Warn("Failed to save config file")
		}
	}()
}

func (u *gtkUI) onActivate() {
	if activeWindow := u.app.GetActiveWindow(); activeWindow != nil {
		activeWindow.Present()
		return
	}

	applyHacks()
	u.mainWindow()

	go u.watchCommands()
	go u.loadConfig(*config.ConfigFile)
}

func (u *gtkUI) Loop() {
	_ = u.app.Connect("activate", u.onActivate)
	u.app.Run([]string{})
}

func (u *gtkUI) initRoster() {
	u.roster = u.newRoster()
}

func (u *gtkUI) mainWindow() {
	builder := newBuilder("Main")
	u.mainBuilder = builder

	builder.ConnectSignals(map[string]interface{}{
		"on_close_window":                       u.quit,
		"on_add_contact_window":                 u.addContactWindow,
		"on_new_conversation":                   u.newCustomConversation,
		"on_about_dialog":                       u.aboutDialog,
		"on_feedback_dialog":                    u.feedbackDialog,
		"on_toggled_check_Item_Merge":           u.toggleMergeAccounts,
		"on_toggled_check_Item_Show_Offline":    u.toggleShowOffline,
		"on_toggled_check_Item_Show_Waiting":    u.toggleShowWaiting,
		"on_toggled_check_Item_Sort_By_Status":  u.toggleSortByStatus,
		"on_toggled_encrypt_configuration_file": u.toggleEncryptedConfig,
		"on_preferences":                        u.showGlobalPreferences,
		"on_muc_show_public_rooms":              u.mucShowPublicRooms,
		"on_muc_show_join_room":                 u.mucShowJoinRoom,
		"on_create_chat_room":                   u.mucShowCreateRoom,
	})

	panicOnDevError(builder.bindObjects(u))

	u.window.SetApplication(u.app)

	u.displaySettings = detectCurrentDisplaySettingsFrom(u.window)
	u.keyboardSettings = newKeyboardSettings()

	// This must happen after u.displaySettings is initialized
	// So now, roster depends on displaySettings which depends on mainWindow
	u.initRoster()

	addItemsThatShouldToggleOnGlobalMenuStatus(builder.getObj("newConvMenu").(isSensitive))
	addItemsThatShouldToggleOnGlobalMenuStatus(builder.getObj("addMenu").(isSensitive))

	// ViewMenu
	u.viewMenu = new(viewMenu)

	panicOnDevError(builder.bindObjects(u.viewMenu))

	u.displaySettings.defaultSettingsOn(u.viewMenu.merge)
	u.displaySettings.defaultSettingsOn(u.viewMenu.offline)
	u.displaySettings.defaultSettingsOn(u.viewMenu.waiting)
	u.displaySettings.defaultSettingsOn(u.viewMenu.sortStatus)

	// OptionsMenu
	u.optionsMenu = new(optionsMenu)
	u.optionsMenu.encryptConfig = builder.getObj("EncryptConfigurationFileCheckMenuItem").(gtki.CheckMenuItem)
	u.displaySettings.defaultSettingsOn(u.optionsMenu.encryptConfig)

	u.initMenuBar()
	obj := builder.getObj("Vbox")
	vbox := obj.(gtki.Box)
	vbox.PackStart(u.roster.widget, true, true, 0)

	obj = builder.getObj("Hbox")
	hbox := obj.(gtki.Box)
	u.unified = newUnifiedLayout(u, vbox, hbox)
	u.unifiedCached = u.unified

	u.config().WhenLoaded(func(a *config.ApplicationConfig) {
		if a.Display.HideFeedbackBar {
			return
		}

		doInUIThread(u.addFeedbackInfoBar)
	})

	u.initSearchBar()

	u.connectShortcutsMainWindow(u.window)

	u.window.SetIcon(coyimIcon.GetPixbuf())
	g.gtk.WindowSetDefaultIcon(coyimIcon.GetPixbuf())

	//Ideally, this should respect widgets initial value for "display",
	//and only call window.Show()
	u.updateGlobalMenuStatus()

	u.initializeMenus()

	u.hooks.BeforeMainWindow(u)

	u.setupSystemTray()

	u.window.ShowAll()
}

func (u *gtkUI) setupSystemTray() {
	si, _ := g.gtk.StatusIconNewFromPixbuf(coyimIcon.GetPixbuf())
	si.SetTooltipText("CoyIM")
	si.SetHasTooltip(true)
	si.SetTitle("CoyIM")
	si.SetVisible(true)
	_ = si.Connect("activate", func() {
		if u.window.IsActive() {
			u.window.Hide()
		} else {
			u.window.Present()
		}
	})
}

func (u *gtkUI) addInitialAccountsToRoster() {
	for _, account := range u.getAllAccounts() {
		u.roster.update(account, rosters.New())
	}
}

func (u *gtkUI) addFeedbackInfoBar() {
	builder := newBuilder("FeedbackInfo")

	obj := builder.getObj("feedbackInfo")
	infobar := obj.(gtki.InfoBar)

	u.notificationArea.PackEnd(infobar, true, true, 0)
	infobar.ShowAll()

	builder.ConnectSignals(map[string]interface{}{
		"handleResponse": func(info gtki.InfoBar, response gtki.ResponseType) {
			if response != gtki.RESPONSE_CLOSE {
				return
			}

			infobar.Hide()
			infobar.Destroy()

			u.config().Display.HideFeedbackBar = true
			u.saveConfigOnly()
		},
	})

	obj = builder.getObj("feedbackButton")
	button := obj.(gtki.Button)
	_ = button.Connect("clicked", func() {
		doInUIThread(u.feedbackDialog)
	})
}

func (u *gtkUI) quit() {
	u.accountManager.disconnectAll()
	u.app.Quit()
}

func (u *gtkUI) askForPassword(accountName string, addGoogleWarning bool, cancel func(), connect func(string) error, savePass func(string)) {
	dialogTemplate := "AskForPassword"

	builder := newBuilder(dialogTemplate)

	dialog := builder.getObj(dialogTemplate).(gtki.Dialog)

	label := builder.getObj("accountName").(gtki.Label)
	label.SetText(accountName)
	label.SetSelectable(true)

	if addGoogleWarning {
		msg := builder.getObj("message").(gtki.Label)
		msg.SetText(i18n.Local("You are trying to connect to a Google account - sometimes Google will not allow connections even if you have entered the correct password. Try turning on App specific password, or if that fails allow less secure applications to access the account (don't worry, CoyIM is plenty secure)."))
		msg.SetSelectable(true)
	}

	passwordEntry := builder.getObj("password").(gtki.Entry)
	savePassword := builder.getObj("savePassword").(gtki.CheckButton)

	builder.ConnectSignals(map[string]interface{}{
		"on_entered_password": func() {
			password, _ := passwordEntry.GetText()
			shouldSave := savePassword.GetActive()

			if len(password) > 0 {
				go func() {
					_ = connect(password)
				}()
				if shouldSave {
					go savePass(password)
				}
				dialog.Destroy()
			}
		},
		"on_cancel_password": func() {
			cancel()
			dialog.Destroy()
		},
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func (u *gtkUI) feedbackDialog() {
	builder := newBuilder("Feedback")

	obj := builder.getObj("dialog")
	dialog := obj.(gtki.Dialog)

	builder.ConnectSignals(map[string]interface{}{
		"on_close": func() {
			dialog.Destroy()
		},
	})

	doInUIThread(func() {
		dialog.SetTransientFor(u.window)
		dialog.ShowAll()
	})
}

func (u *gtkUI) shouldViewAccounts() bool {
	return !u.config().Display.MergeAccounts
}

func (u *gtkUI) aboutDialog() {
	//TODO: This dialog automatically parses HTML and display clickable links.
	//We may need to use  g_markup_escape_text().
	dialog, _ := g.gtk.AboutDialogNew()
	dialog.SetName(i18n.Local("CoyIM!"))
	dialog.SetProgramName(programName)
	dialog.SetAuthors(authors())
	dialog.SetVersion(coyimVersion)
	dialog.SetLicense(`GNU GENERAL PUBLIC LICENSE, Version 3`)
	dialog.SetWrapLicense(true)

	dialog.SetTransientFor(u.window)
	dialog.Run()
	dialog.Destroy()
}

func (u *gtkUI) newCustomConversation() {
	accounts := u.getAllConnectedAccounts()

	var dialog gtki.Window
	var model gtki.ListStore
	var accountInput gtki.ComboBox
	var peerInput gtki.Entry

	builder := newBuilder("NewCustomConversation")
	builder.getItems(
		"NewCustomConversation", &dialog,
		"accounts-model", &model,
		"accounts", &accountInput,
		"address", &peerInput,
	)

	dialog.SetApplication(u.app)

	for _, acc := range accounts {
		iter := model.Append()
		_ = model.SetValue(iter, 0, acc.Account())
		_ = model.SetValue(iter, 1, acc.ID())
	}

	if len(accounts) > 0 {
		accountInput.SetActive(0)
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_close": dialog.Destroy,
		"on_start": func() {
			iter, err := accountInput.GetActiveIter()
			if err != nil {
				u.log.WithError(err).Warn("Error encountered when getting account")
				return
			}
			val, err := model.GetValue(iter, 1)
			if err != nil {
				u.log.WithError(err).Warn("Error encountered when getting account")
				return
			}
			accountID, _ := val.GetString()

			account, ok := u.accountManager.getAccountByID(accountID)
			if !ok {
				return
			}
			j, _ := peerInput.GetText()
			jj := jid.Parse(j)
			switch jjj := jj.(type) {
			case jid.WithResource:
				u.openTargetedConversationView(account, jjj, true)
			default:
				u.openConversationView(account, jj, true)
			}

			dialog.Destroy()
		},
	})

	dialog.SetTransientFor(u.window)
	dialog.ShowAll()
}

func (u *gtkUI) addContactWindow() {
	dialog := u.presenceSubscriptionDialog(func(accountID string, peer jid.WithoutResource, msg, nick string, autoAuth bool) error {
		account, ok := u.accountManager.getAccountByID(accountID)
		if !ok {
			return fmt.Errorf(i18n.Local("There is no account with the id %q"), accountID)
		}

		if !account.connected() {
			return errors.New(i18n.Local("Can't send a contact request from an offline account"))
		}

		err := account.session.RequestPresenceSubscription(peer, msg)
		rl := u.accountManager.getContacts(account)
		rl.SubscribeRequest(peer, "", accountID)

		if nick != "" {
			account.session.GetConfig().SavePeerDetails(peer.String(), nick, []string{})
			u.SaveConfig()
		}

		if autoAuth {
			account.session.AutoApprove(peer.String())
		}

		return err
	})

	dialog.SetTransientFor(u.window)
	dialog.Show()
}

func (u *gtkUI) listenToToggleConnectAllAutomatically() {
	for {
		val := <-u.toggleConnectAllAutomaticallyRequest
		u.config().ConnectAutomatically = val
		u.saveConfigOnly()
	}
}

func (u *gtkUI) setConnectAllAutomatically(val bool) {
	u.toggleConnectAllAutomaticallyRequest <- val
}

func (u *gtkUI) setShowAdvancedSettings(val bool) {
	u.setShowAdvancedSettingsRequest <- val
}

func (u *gtkUI) listenToSetShowAdvancedSettings() {
	for {
		val := <-u.setShowAdvancedSettingsRequest
		u.config().AdvancedOptions = val
		u.saveConfigOnly()
	}
}

func (u *gtkUI) setMenuBarSensitive(v bool) {
	u.contactsMenuItem.SetSensitive(v)
	u.accountsMenuItem.SetSensitive(v)
	u.chatRoomsMenuItem.SetSensitive(v)
	u.viewMenuItem.SetSensitive(v)
	u.optionsMenuItem.SetSensitive(v)
}

func (u *gtkUI) initMenuBar() {
	_ = u.window.Connect(accountChangedSignal.String(), func() {
		doInUIThread(func() {
			u.buildAccountsMenu()
			u.accountsMenuItem.ShowAll()
			u.rosterUpdated()
		})
	})

	u.buildAccountsMenu()
	u.accountsMenuItem.ShowAll()

	u.setMenuBarSensitive(false)

	u.whenHaveConfig(func() {
		doInUIThread(func() {
			u.setMenuBarSensitive(true)
		})
	})
}

func (u *gtkUI) initSearchBar() {
	u.searchEntry.SetCanFocus(true)
	u.searchEntry.Map()

	u.search.SetHAlign(gtki.ALIGN_FILL)
	u.search.SetHExpand(true)
	u.search.ConnectEntry(u.searchEntry)
	u.roster.view.SetSearchEntry(u.searchEntry)

	prov := providerWithCSS("entry { min-width: 300px; }")
	updateWithStyle(u.searchEntry, prov)

	prov = providerWithCSS("box { border: none; }")
	updateWithStyle(u.searchBox, prov)

	// TODO: unify with dark themes
	prov = providerWithCSS("searchbar {background-color: #e8e8e7; }")
	updateWithStyle(u.search, prov)
}

func (u *gtkUI) rosterUpdated() {
	doInUIThread(u.roster.redraw)
	if u.unified != nil {
		doInUIThread(u.unified.update)
	}
}

func (u *gtkUI) editAccount(account *account) {
	u.accountDialog(account.session, account.session.GetConfig(), func() {
		u.SaveConfig()
		account.session.ReloadKeys()
	})
}

func (u *gtkUI) removeAccount(account *account) {
	u.confirmAccountRemoval(account.session.GetConfig(), func(c *config.Account) {
		account.session.SetWantToBeOnline(false)
		account.disconnect()
		u.removeSaveReload(c)
	})
}

func (u *gtkUI) toggleAutoConnectAccount(account *account) {
	account.session.GetConfig().ToggleConnectAutomatically()
	u.saveConfigOnly()
}

func (u *gtkUI) presenceUpdated(account *account, peer jid.WithResource, ev events.Presence) {
	//TODO: Ignore presence updates for yourself.
	if account == nil {
		return
	}

	u.NewConversationViewFactory(account, peer, false).IfConversationView(func(c conversationView) {
		doInUIThread(func() {
			c.appendStatus(u.roster.displayNameFor(account, peer.NoResource()), time.Now(), ev.Show, ev.Status, ev.Gone)
		})
	}, func() {})
}

func (u *gtkUI) toggleAlwaysEncryptAccount(account *account) {
	account.session.GetConfig().ToggleAlwaysEncrypt()
	u.saveConfigOnly()
}

func (u *gtkUI) openConversationView(account *account, peer jid.Any, userInitiated bool) conversationView {
	return u.NewConversationViewFactory(account, peer, false).OpenConversationView(userInitiated)
}

func (u *gtkUI) openTargetedConversationView(account *account, peer jid.Any, userInitiated bool) conversationView {
	return u.NewConversationViewFactory(account, peer, true).OpenConversationView(userInitiated)
}

func (u *gtkUI) config() *config.ApplicationConfig {
	u.internalConfigLock.Lock()
	defer u.internalConfigLock.Unlock()

	c := u.internalConfig

	return c
}

func (u *gtkUI) setConfig(c *config.ApplicationConfig) {
	u.internalConfigLock.Lock()
	u.internalConfig = c
	u.internalConfigLock.Unlock()

	u.haveConfig()
}

func (u *gtkUI) whenHaveConfig(f func()) {
	if u.config() != nil {
		f()
		return
	}
	u.haveConfigEntries.add(f)
}

func (u *gtkUI) haveConfig() {
	cs := u.haveConfigEntries
	u.haveConfigEntries = nil
	cs.invokeAll()
}
