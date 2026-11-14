package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"
	"github.com/xmhrnx/giode/components/titlebar"

	"gopdf/ui/components"
	"gopdf/ui/platform"
	"gopdf/ui/theme"
)

// titleMaxWidth caps the document name in the title bar, in dp.
const titleMaxWidth = 320

// titleBar is the custom window frame: logo, a menu bar that opens on
// hover, the document name centered after it, and window buttons.
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
	mb := giode.MenuBar().OpenOnHover(true)
	// Each panel is laid out beside its trigger: laid out elsewhere, a
	// panel opened by hovering shows under the first menu instead.
	item := func(menu *menu.Menu, label string) giode.Element {
		return giode.Row(giode.Styles{}, mb.Add(menu, components.MenuLabel(label)), menu)
	}
	t.bar = giode.TitleBar(s.win, "GoPDF", giode.Styles{
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
		TitlePlacement(giode.TitleBarTitleCenter).
		TitleMaxWidth(titleMaxWidth).
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
