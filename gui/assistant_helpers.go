package gui

import "github.com/coyim/gotk3adapter/gtki"

// findAssistantHeaderContainer MUST be called from the UI thread
func findAssistantHeaderContainer(a gtki.Assistant) gtki.Container {
	lbl, _ := g.gtk.LabelNew("")
	a.AddActionWidget(lbl)

	parentBox, _ := lbl.GetParentX()
	a.RemoveActionWidget(lbl)

	return parentBox.(gtki.Container)
}

const (
	assistantButtonBackCancelName = "cancel"
	assistantButtonBackLastName   = "back"
	assistantButtonLastName       = "last"
	assistantButtonForwardName    = "forward"
	assistantButtonApplyName      = "apply"
)

var assistantNavigationButtons = []string{
	assistantButtonBackLastName,
	assistantButtonLastName,
	assistantButtonForwardName,
	assistantButtonApplyName,
}

type assistantButtons map[string]gtki.Button

// getButtonsForAssistantHeader MUST be called from the UI thread
func getButtonsForAssistantHeader(a gtki.Assistant) assistantButtons {
	h := findAssistantHeaderContainer(a)
	result := assistantButtons{}

	for _, c := range h.GetChildren() {
		if b, ok := c.(gtki.Button); ok {
			name, _ := g.gtk.GetWidgetBuildableName(b)
			result[name] = b
		}
	}

	return result
}

// updateLastButtonLabel MUST be called from the UI thread
func (list assistantButtons) updateLastButtonLabel(label string) {
	list.updateButtonLabelByName(assistantButtonLastName, label)
}

// updateApplyButtonLabel MUST be called from the UI thread
func (list assistantButtons) updateApplyButtonLabel(label string) {
	list.updateButtonLabelByName(assistantButtonApplyName, label)
}

// updateButtonLabelByName MUST be called from the UI thread
func (list assistantButtons) updateButtonLabelByName(name string, label string) {
	if b, ok := list[name]; ok {
		b.SetLabel(label)
	}
}

// disableNavigationButNotCancel MUST be called from the UI thread
func (list assistantButtons) disableNavigationButNotCancel() {
	for _, buttonName := range assistantNavigationButtons {
		if b, ok := list[buttonName]; ok {
			b.SetSensitive(false)
		}
	}
}

// enableNavigation MUST be called from the UI thread
func (list assistantButtons) enableNavigation() {
	for _, buttonName := range assistantNavigationButtons {
		if b, ok := list[buttonName]; ok {
			b.SetSensitive(true)
		}
	}
}

const (
	assistantActionAreaName     = "action_area"
	assistantSidebarName        = "sidebar"
	assistantContentWrapperName = "content_box"
	assistantContentName        = "content"
)

// getBottomActionAreaFromAssistant MUST be called from the UI thread
func getBottomActionAreaFromAssistant(a gtki.Assistant) (gtki.Box, bool) {
	return findGtkBoxWithID(a.GetChildren(), assistantActionAreaName)
}

// getSidebarFromAssistant MUST be called from the UI thread
func getSidebarFromAssistant(a gtki.Assistant) (gtki.Box, bool) {
	return findGtkBoxWithID(a.GetChildren(), assistantSidebarName)
}

// getPagesFromAssistant MUST be called from the UI thread
func getPagesFromAssistant(a gtki.Assistant) []gtki.Widget {
	if notebook, ok := getNotebookFromAssistant(a); ok {
		return notebook.GetChildren()
	}
	return nil
}

// getNotebookFromAssistant MUST be called from the UI thread
func getNotebookFromAssistant(a gtki.Assistant) (gtki.Notebook, bool) {
	if content, ok := findGtkBoxWithID(a.GetChildren(), assistantContentWrapperName); ok {
		for _, ch := range content.GetChildren() {
			if notebook, ok := ch.(gtki.Notebook); ok {
				if name, _ := g.gtk.GetWidgetBuildableName(notebook); name == assistantContentName {
					return notebook, true
				}
			}
		}
	}
	return nil, false
}

// removeMarginFromAssistantPages MUST be called from the UI thread
func removeMarginFromAssistantPages(a gtki.Assistant) {
	for _, page := range getPagesFromAssistant(a) {
		page.SetProperty("margin", 0)
	}
}

// setAssistantSidebar MUST be called from the UI thread
func setAssistantSidebarContent(a gtki.Assistant, content gtki.Widget) {
	if sidebar, ok := getSidebarFromAssistant(a); ok {
		for _, ch := range sidebar.GetChildren() {
			sidebar.Remove(ch)
		}
		sidebar.PackStart(content, false, false, 0)
	}
}

// findGtkBoxWithID MUST be called from the UI thread
func findGtkBoxWithID(list []gtki.Widget, boxName string) (gtki.Box, bool) {
	for _, widget := range list {
		if box, ok := widget.(gtki.Box); ok {
			if name, _ := g.gtk.GetWidgetBuildableName(box); name == boxName {
				return box, true
			}
			if box, ok = findGtkBoxWithID(box.GetChildren(), boxName); ok {
				return box, true
			}
		}
	}
	return nil, false
}
