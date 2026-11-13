package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

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
		giode.Shortcut("Delete", s.ui.unlessTyping(s.deleteSelected)),
		giode.Shortcut("Backspace", s.ui.unlessTyping(s.deleteSelected)),
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
		giode.Shortcut("PageDown", s.ui.unlessTyping(func() { s.goToPage(s.view.current + 1) })),
		giode.Shortcut("PageUp", s.ui.unlessTyping(func() { s.goToPage(s.view.current - 1) })),
		giode.Shortcut("Mod+R", cur(func(p int) { s.rotatePage(p, true) })),
		giode.Shortcut("Mod+Shift+R", cur(func(p int) { s.rotatePage(p, false) })),
		giode.Shortcut("Mod+D", cur(s.duplicatePage)),
		giode.Shortcut("Mod+Backspace", cur(s.deletePage)),
	}
	if s.hasDoc() {
		for _, d := range tools {
			tool := d.tool
			list = append(list, giode.Shortcut(d.key, s.ui.unlessTyping(func() { s.setTool(tool) })))
		}
	}
	giode.Shortcuts(list...)
}
