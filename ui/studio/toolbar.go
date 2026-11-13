package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// toolbar is the strip under the title bar: file actions, history,
// drawing tools and zoom.
type toolbar struct {
	s       *Studio
	sidebar *components.ToolButton
	open    *components.ToolButton
	save    *components.ToolButton
	undo    *components.ToolButton
	redo    *components.ToolButton
	tools   []*components.ToolButton
	zoom    *zoomControl
}

func newToolbar(s *Studio) *toolbar {
	t := &toolbar{s: s, zoom: newZoomControl(s)}
	t.sidebar = components.NewToolButton("panel-left", "Toggle sidebar  ("+giode.FormatShortcut("Mod+B")+")", s.toggleSidebar).
		Active(func() bool { return s.view.sidebar })
	t.open = components.NewToolButton("folder-open", "Open  ("+giode.FormatShortcut("Mod+O")+")", s.openDialog)
	t.save = components.NewToolButton("save", "Save  ("+giode.FormatShortcut("Mod+S")+")", s.save).
		Enabled(func() bool { return s.canEdit() && s.doc.Dirty() })
	t.undo = components.NewToolButton("undo-2", "Undo  ("+giode.FormatShortcut("Mod+Z")+")", s.undo).
		Enabled(func() bool { return s.canEdit() && s.doc.CanUndo() })
	t.redo = components.NewToolButton("redo-2", "Redo  ("+giode.FormatShortcut("Mod+Shift+Z")+")", s.redo).
		Enabled(func() bool { return s.canEdit() && s.doc.CanRedo() })
	for _, d := range tools {
		tool := d.tool
		b := components.NewToolButton(d.icon, d.name+"  ("+d.key+")", func() { s.setTool(tool) }).
			Active(func() bool { return s.tool == tool })
		t.tools = append(t.tools, b)
	}
	return t
}

// View returns the toolbar for this frame.
func (t *toolbar) View() giode.Element {
	left := giode.Row(giode.Styles{Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		t.sidebar.View(), components.VDivider(20),
		t.open.View(), t.save.View(), components.VDivider(20),
		t.undo.View(), t.redo.View(),
	)
	toolViews := make([]giode.Element, len(t.tools))
	for i, b := range t.tools {
		toolViews[i] = b.View()
	}
	center := giode.Row(giode.Styles{
		Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.UniformInset(3),
		Background: theme.Subtle, BorderRadius: 11, BorderWidth: 1, BorderColor: theme.Border,
	}, toolViews...)
	bar := giode.Row(giode.Styles{
		Width: -1, Height: theme.ToolbarHeight, Align: giode.AlignCenter, Justify: giode.JustifySpaceBetween,
		Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(10, 0), Background: theme.Surface,
	}, left, center, t.zoom.View())
	return giode.Stack(giode.Styles{Width: -1}, bar, components.HRule())
}
