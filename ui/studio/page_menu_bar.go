package studio

import (
	"github.com/xmhrnx/giode/components/menu"

	"gopdf/ui/theme"
)

// newPageMenuBar is the title bar's Page menu; it acts on the current page.
func newPageMenuBar(s *Studio) *menu.Menu {
	cur := func(fn func(int)) func() { return func() { fn(s.view.current) } }
	rotateAll := newMenu().
		Item("Clockwise", func() { s.rotateAll(true) }).
		Item("Counter-clockwise", func() { s.rotateAll(false) })
	return newMenu().
		IconItem(theme.MenuIcon("rotate-cw"), "Rotate Clockwise", cur(func(p int) { s.rotatePage(p, true) })).Shortcut("Mod+R").
		IconItem(theme.MenuIcon("rotate-ccw"), "Rotate Counter-clockwise", cur(func(p int) { s.rotatePage(p, false) })).Shortcut("Mod+Shift+R").
		SubMenu("Rotate All Pages", rotateAll).
		Separator().
		IconItem(theme.MenuIcon("copy"), "Duplicate Page", cur(s.duplicatePage)).
		IconItem(theme.MenuIcon("file-plus-2"), "Insert Blank Page After", cur(s.insertBlankPage)).
		IconItem(theme.MenuIcon("arrow-up"), "Move Page Up", cur(func(p int) { s.movePage(p, p-1) })).
		IconItem(theme.MenuIcon("arrow-down"), "Move Page Down", cur(func(p int) { s.movePage(p, p+1) })).
		Separator().
		IconItem(theme.MenuIcon("trash-2"), "Delete Page", cur(s.deletePage)).Shortcut("Mod+Backspace")
}
