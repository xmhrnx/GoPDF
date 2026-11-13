package studio

import (
	"github.com/xmhrnx/giode/components/menu"

	"gopdf/ui/theme"
)

func newViewMenu(s *Studio) *menu.Menu {
	return newMenu().
		IconItem(theme.MenuIcon("zoom-in"), "Zoom In", func() { s.zoomBy(1) }).Shortcut("Mod+=").
		IconItem(theme.MenuIcon("zoom-out"), "Zoom Out", func() { s.zoomBy(-1) }).Shortcut("Mod+-").
		IconItem(theme.MenuIcon("move-horizontal"), "Fit Width", s.fitWidth).Shortcut("Mod+0").
		Item("Actual Size", func() { s.setZoom(1) }).Shortcut("Mod+1").
		Separator().
		CheckItem("Sidebar", func() bool { return s.view.sidebar }, s.toggleSidebar).Shortcut("Mod+B").
		CheckItem("Inspector", func() bool { return s.view.inspector }, s.toggleInspector).Shortcut("Mod+I").
		Separator().
		CheckItem("Dark Mode", theme.IsDark, theme.ToggleMode).Shortcut("Mod+Shift+D")
}
