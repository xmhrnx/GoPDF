package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// views holds every stateful view, created once.
type views struct {
	s *Studio

	titleBar  *titleBar
	toolbar   *toolbar
	sidebar   *sidebar
	search    *searchPanel
	viewer    *viewer
	inspector *inspector
	status    *statusBar
	welcome   *welcomeScreen
	pageMenu  *pageMenu

	password    *passwordDialog
	unsaved     *unsavedDialog
	properties  *propertiesDialog
	watermark   *watermarkDialog
	pageNumbers *pageNumbersDialog
	security    *securityDialog
	extract     *extractDialog
	split       *splitDialog
	shortcuts   *shortcutsDialog
	signer      *signatureDialog
}

func newViews(s *Studio) *views {
	v := &views{s: s}
	s.ui = v
	v.pageMenu = newPageMenu(s)
	v.search = newSearchPanel(s)
	v.titleBar = newTitleBar(s)
	v.toolbar = newToolbar(s)
	v.sidebar = newSidebar(s)
	v.viewer = newViewer(s)
	v.inspector = newInspector(s)
	v.status = newStatusBar(s)
	v.welcome = newWelcomeScreen(s)
	v.password = newPasswordDialog(s)
	v.unsaved = newUnsavedDialog(s)
	v.properties = newPropertiesDialog(s)
	v.watermark = newWatermarkDialog(s)
	v.pageNumbers = newPageNumbersDialog(s)
	v.security = newSecurityDialog(s)
	v.extract = newExtractDialog(s)
	v.split = newSplitDialog(s)
	v.shortcuts = newShortcutsDialog(s)
	v.signer = newSignatureDialog(s)
	return v
}

// root is the whole window.
func (v *views) root() giode.Element {
	registerShortcuts(v.s)
	var body giode.Element
	if v.s.hasDoc() {
		body = v.workspace()
	} else {
		body = v.welcome.View()
	}
	kids := append([]giode.Element{v.titleBar.View(), body}, v.overlays()...)
	return giode.Stack(giode.Styles{Width: -1, Height: -1, Background: theme.Background}, kids...)
}

func (v *views) workspace() giode.Element {
	row := []giode.Element{}
	if v.s.view.sidebar {
		row = append(row, v.sidebar.View())
	}
	row = append(row, v.viewer.View())
	if v.s.view.inspector {
		row = append(row, v.inspector.View())
	}
	return giode.Stack(giode.Styles{Width: -1, FlexGrow: 1},
		v.toolbar.View(),
		giode.Row(giode.Styles{Width: -1, FlexGrow: 1, Wrap: giode.WrapNoWrap, Align: giode.AlignStretch}, row...),
		v.status.View(),
	)
}

// overlays paint above the window and take no layout space.
func (v *views) overlays() []giode.Element {
	return []giode.Element{
		dropOverlay(v.s.dragOver),
		v.pageMenu.menu,
		v.password.dialog, v.unsaved.dialog, v.properties.dialog, v.watermark.dialog,
		v.pageNumbers.dialog, v.security.dialog, v.extract.dialog, v.split.dialog, v.shortcuts.dialog, v.signer.dialog,
		v.s.toasts,
	}
}
