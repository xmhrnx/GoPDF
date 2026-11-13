package studio

import (
	"github.com/xmhrnx/giode/components/menu"

	"gopdf/ui/theme"
)

func newFileMenu(s *Studio) *menu.Menu {
	return newMenu().
		IconItem(theme.MenuIcon("folder-open"), "Open…", s.openDialog).Shortcut("Mod+O").
		IconItem(theme.MenuIcon("save"), "Save", s.save).Shortcut("Mod+S").
		IconItem(theme.MenuIcon("save-all"), "Save As…", s.saveAs).Shortcut("Mod+Shift+S").
		Separator().
		IconItem(theme.MenuIcon("files"), "Merge PDFs…", s.mergeDialog).
		IconItem(theme.MenuIcon("file-output"), "Extract Pages…", s.ui.openExtract).
		IconItem(theme.MenuIcon("scissors"), "Split Document…", s.ui.openSplit).
		Separator().
		IconItem(theme.MenuIcon("image"), "Export Page as PNG…", s.exportPageImage).
		IconItem(theme.MenuIcon("file-text"), "Export Text…", s.exportText).
		Separator().
		IconItem(theme.MenuIcon("file-x"), "Close Document", s.closeDocument).Shortcut("Mod+W")
}
