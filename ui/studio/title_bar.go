package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/titlebar"

	"gopdf/ui/components"
	"gopdf/ui/platform"
	"gopdf/ui/theme"
)

// titleBar is the custom window frame: logo, menus, document name and
// window buttons.
type titleBar struct {
	s     *Studio
	bar   *titlebar.TitleBar
	theme *components.ToolButton
	panel *components.ToolButton
}

func newTitleBar(s *Studio) *titleBar {
	t := &titleBar{s: s}
	t.theme = components.NewToolButton("moon", "Toggle dark mode", theme.ToggleMode)
	t.panel = components.NewToolButton("panel-right", "Toggle inspector", s.toggleInspector).
		Active(func() bool { return s.view.inspector }).Enabled(s.hasDoc)
	m := newMenus(s)
	t.bar = giode.TitleBar(s.win, "GoPDF", giode.Styles{
		Height: theme.TitleBarHeight, Background: theme.Surface, Color: theme.TextColor,
		BorderWidth: 1, BorderColor: theme.Border, FontSize: 12.5,
	}).
		Variant(platform.TitleBarVariant).
		Icon(appLogo(s.logo, 24)).
		Leading(
			components.NewMenuButton("File", m.file),
			components.NewMenuButton("Edit", m.edit),
			components.NewMenuButton("View", m.view),
			components.NewMenuButton("Page", m.page),
			components.NewMenuButton("Tools", m.tools),
			components.NewMenuButton("Help", m.help),
		).
		Trailing(giode.Raw(t.trailing)).
		OnClose(s.OnClose)
	return t
}

// View returns the bar with the current document name.
func (t *titleBar) View() giode.Element {
	t.bar.Title(t.title())
	return t.bar
}

func (t *titleBar) title() string {
	if !t.s.hasDoc() {
		return "GoPDF"
	}
	if t.s.doc.Dirty() {
		return t.s.doc.Name + "  •  Edited"
	}
	return t.s.doc.Name
}

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
