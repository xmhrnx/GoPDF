package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// nudgeKeys move the selection with the arrow keys; Shift moves further.
var nudgeKeys = []struct {
	spec   string
	dx, dy float64
	large  bool
}{
	{"Left", -1, 0, false}, {"Right", 1, 0, false}, {"Up", 0, -1, false}, {"Down", 0, 1, false},
	{"Shift+Left", -1, 0, true}, {"Shift+Right", 1, 0, true}, {"Shift+Up", 0, -1, true}, {"Shift+Down", 0, 1, true},
}

// registerShortcuts installs the keyboard shortcuts for this frame.
func registerShortcuts(s *Studio) {
	cur := func(fn func(int)) func() {
		return func() {
			if s.hasDoc() {
				fn(s.view.current)
			}
		}
	}
	list := []giode.KeyShortcut{
		giode.Shortcut("Mod+O", s.openDialog),
		giode.Shortcut("Mod+S", s.save),
		giode.Shortcut("Mod+Shift+S", s.saveAs),
		giode.Shortcut("Mod+W", s.closeDocument),
		giode.Shortcut("Mod+Z", s.undo),
		giode.Shortcut("Mod+Shift+Z", s.redo),
		giode.Shortcut("Mod+Y", s.redo),
		giode.Shortcut("Delete", s.ui.unlessDialog(s.deleteSelected)),
		giode.Shortcut("Backspace", s.ui.unlessDialog(s.deleteSelected)),
		giode.Shortcut("Esc", s.escapeKey),
		giode.Shortcut("Mod+=", func() { s.zoomBy(1) }),
		giode.Shortcut("Mod+Shift+=", func() { s.zoomBy(1) }),
		giode.Shortcut("Mod+-", func() { s.zoomBy(-1) }),
		giode.Shortcut("Mod+0", s.fitWidth),
		giode.Shortcut("Mod+1", func() { s.setZoom(1) }),
		giode.Shortcut("Mod+B", s.toggleSidebar),
		giode.Shortcut("Mod+I", s.toggleInspector),
		giode.Shortcut("Mod+F", s.ui.search.focus),
		giode.Shortcut("Mod+/", s.ui.openShortcuts),
		giode.Shortcut("Mod+Shift+D", theme.ToggleMode),
		giode.Shortcut("PageDown", s.ui.unlessDialog(func() { s.goToPage(s.view.current + 1) })),
		giode.Shortcut("PageUp", s.ui.unlessDialog(func() { s.goToPage(s.view.current - 1) })),
		giode.Shortcut("Mod+R", cur(func(p int) { s.rotatePage(p, true) })),
		giode.Shortcut("Mod+Shift+R", cur(func(p int) { s.rotatePage(p, false) })),
		giode.Shortcut("Mod+D", s.duplicate),
		giode.Shortcut("Mod+C", s.ui.unlessEditing(s.copySelection)),
		giode.Shortcut("Mod+X", s.ui.unlessEditing(s.cutSelection)),
		giode.Shortcut("Mod+V", s.ui.unlessEditing(s.paste)),
		giode.Shortcut("Mod+Backspace", cur(s.deletePage)),
	}
	for _, n := range nudgeKeys {
		list = append(list, giode.Shortcut(n.spec, s.ui.unlessDialog(func() { s.nudge(n.dx, n.dy, n.large) })))
	}
	if s.hasDoc() {
		for _, d := range tools {
			tool := d.tool
			list = append(list, giode.Shortcut(d.key, s.ui.unlessDialog(func() { s.setTool(tool) })))
		}
	}
	giode.Shortcuts(list...)
}
