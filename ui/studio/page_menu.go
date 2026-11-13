package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"

	"gopdf/ui/theme"
)

// pageMenu is the context menu of a page thumbnail.
type pageMenu struct {
	s    *Studio
	page int
	menu *menu.Menu
}

func newPageMenu(s *Studio) *pageMenu {
	m := &pageMenu{s: s}
	at := func(fn func(int)) func() { return func() { fn(m.page) } }
	m.menu = newMenu().
		IconItem(theme.MenuIcon("rotate-cw"), "Rotate Clockwise", at(func(p int) { s.rotatePage(p, true) })).
		IconItem(theme.MenuIcon("rotate-ccw"), "Rotate Counter-clockwise", at(func(p int) { s.rotatePage(p, false) })).
		Separator().
		IconItem(theme.MenuIcon("copy"), "Duplicate", at(s.duplicatePage)).
		IconItem(theme.MenuIcon("file-plus-2"), "Insert Blank Page After", at(s.insertBlankPage)).
		IconItem(theme.MenuIcon("arrow-up"), "Move Up", at(func(p int) { s.movePage(p, p-1) })).
		IconItem(theme.MenuIcon("arrow-down"), "Move Down", at(func(p int) { s.movePage(p, p+1) })).
		IconItem(theme.MenuIcon("file-output"), "Extract This Page…", at(s.ui.extractPage)).
		Separator().
		IconItem(theme.MenuIcon("trash-2"), "Delete Page", at(s.deletePage))
	return m
}

// open shows the menu for page where e happened.
func (m *pageMenu) open(page int, e giode.ClickEvent) {
	m.page = page
	m.menu.OpenAtClick(e)
}
