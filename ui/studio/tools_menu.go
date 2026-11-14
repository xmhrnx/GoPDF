package studio

import (
	"github.com/xmhrnx/giode/components/menu"

	"gopdf/ui/theme"
)

func newToolsMenu(s *Studio) *menu.Menu {
	return newMenu().
		IconItem(theme.MenuIcon("stamp"), "Add Watermark…", s.ui.openWatermark).
		IconItem(theme.MenuIcon("hash"), "Add Page Numbers…", s.ui.openPageNumbers).
		Separator().
		IconItem(theme.MenuIcon("lock"), "Password Protection…", s.ui.openSecurity).
		IconItem(theme.MenuIcon("badge-check"), "Digital IDs…", func() { s.ui.digitalIDs.show(nil) }).
		IconItem(theme.MenuIcon("shield-check"), "Signature Details…", s.ui.openSignatureDetails).
		IconItem(theme.MenuIcon("info"), "Document Properties…", s.ui.openProperties)
}
