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
		IconItem(theme.MenuIcon("scissors"), "Cut", s.cutSelection).Shortcut("Mod+X").
		IconItem(theme.MenuIcon("clipboard-copy"), "Copy", s.copySelection).Shortcut("Mod+C").
		IconItem(theme.MenuIcon("clipboard-paste"), "Paste", s.paste).Shortcut("Mod+V").
		IconItem(theme.MenuIcon("copy"), "Duplicate", s.duplicate).Shortcut("Mod+D").
		IconItem(theme.MenuIcon("trash-2"), "Delete Selection", s.deleteSelected).Shortcut("Delete").
		Separator().
		IconItem(theme.MenuIcon("signature"), "New Signature…", func() { s.ui.signaturePad.show() }).
		Separator().
		IconItem(theme.MenuIcon("info"), "Document Properties…", s.ui.openProperties)
}
