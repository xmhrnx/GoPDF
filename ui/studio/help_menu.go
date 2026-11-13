package studio

import (
	"github.com/xmhrnx/giode/components/menu"

	"gopdf/ui/theme"
)

func newHelpMenu(s *Studio) *menu.Menu {
	return newMenu().
		IconItem(theme.MenuIcon("keyboard"), "Keyboard Shortcuts", s.ui.openShortcuts).Shortcut("Mod+/")
}
