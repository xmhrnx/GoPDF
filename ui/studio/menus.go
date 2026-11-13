package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"

	"gopdf/ui/theme"
)

// menus are the title bar's drop-down menus.
type menus struct {
	file, edit, view, page, tools, help *menu.Menu
}

// menuStyles is the look of every menu panel.
var menuStyles = giode.Styles{Background: theme.Surface, Color: theme.TextColor, BorderRadius: 10}

// newMenu returns an empty menu in the app's style.
func newMenu() *menu.Menu {
	return giode.Menu(menuStyles).SeparatorColor(theme.Border).HoverColor(theme.Hover)
}

func newMenus(s *Studio) *menus {
	return &menus{
		file:  newFileMenu(s),
		edit:  newEditMenu(s),
		view:  newViewMenu(s),
		page:  newPageMenuBar(s),
		tools: newToolsMenu(s),
		help:  newHelpMenu(s),
	}
}
