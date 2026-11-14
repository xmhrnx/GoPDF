package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"
	"github.com/xmhrnx/giode/components/titlebar"

	"gopdf/ui/components"
	"gopdf/ui/platform"
	"gopdf/ui/theme"
)

// titleBar is the custom window frame: logo, a menu bar that opens on
// hover, and window buttons. Document names belong to the tabs.
type titleBar struct {
	s     *Studio
	bar   *titlebar.TitleBar
	menus *menus
	theme *components.ToolButton
	panel *components.ToolButton
}

func newTitleBar(s *Studio) *titleBar {
	t := &titleBar{s: s}
	t.theme = components.NewToolButton("moon", "Toggle dark mode", theme.ToggleMode)
	t.panel = components.NewToolButton("panel-right", "Toggle inspector", s.toggleInspector).
		Active(func() bool { return s.view.inspector }).Enabled(s.hasDoc)
	m := newMenus(s)
	mb := giode.MenuBar().OpenOnHover(true)
	item := func(menu *menu.Menu, label string) giode.Element {
		return mb.Add(menu, components.MenuLabel(label))
	}
	t.menus = m
	t.bar = giode.TitleBar(s.win, "", giode.Styles{
		Height: theme.TitleBarHeight, Background: theme.Surface, Color: theme.TextColor,
		BorderWidth: 1, BorderColor: theme.Border, FontSize: 12.5,
	}).
		Variant(platform.TitleBarVariant).
		Icon(appLogo(s.logo, 24)).
		Leading(
			item(m.file, "File"),
			item(m.edit, "Edit"),
			item(m.view, "View"),
			item(m.page, "Page"),
			item(m.tools, "Tools"),
			item(m.help, "Help"),
		).
		Trailing(giode.Raw(t.trailing)).
		OnClose(s.OnClose)
	return t
}

// panels are the menus' drop-down panels, laid out with the overlays.
func (t *titleBar) panels() []giode.Element {
	m := t.menus
	return []giode.Element{m.file, m.edit, m.view, m.page, m.tools, m.help}
}

// View returns the bar for this frame.
func (t *titleBar) View() giode.Element { return t.bar }

func (t *titleBar) trailing(gtx giode.Context) giode.Dimensions {
	icon := "moon"
	if theme.IsDark() {
		icon = "sun"
	}
	t.theme.SetIcon(icon)
	return giode.Row(giode.Styles{Gap: 2, Align: giode.AlignCenter, Padding: giode.Inset{Right: 6}},
		t.theme.View(), t.panel.View(),
	).Layout(gtx)
}
