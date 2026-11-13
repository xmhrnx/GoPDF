package studio

import (
	"github.com/xmhrnx/giode/components/menu"

	"gopdf/ui/theme"
)

func newEditMenu(s *Studio) *menu.Menu {
	return newMenu().
		IconItem(theme.MenuIcon("undo-2"), "Undo", s.undo).Shortcut("Mod+Z").
		IconItem(theme.MenuIcon("redo-2"), "Redo", s.redo).Shortcut("Mod+Shift+Z").
		Separator().
		IconItem(theme.MenuIcon("trash-2"), "Delete Selection", s.deleteSelected).Shortcut("Delete").
		Separator().
		IconItem(theme.MenuIcon("info"), "Document Properties…", s.ui.openProperties)
}
