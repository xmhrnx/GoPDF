package studio

// shortcutHelp is one line of the keyboard shortcuts dialog.
type shortcutHelp struct {
	spec, what string
}

// shortcutGroups lists the shortcuts by topic, as the dialog shows them.
var shortcutGroups = []struct {
	title string
	items []shortcutHelp
}{
	{"File", []shortcutHelp{
		{"Mod+O", "Open"}, {"Mod+S", "Save"}, {"Mod+Shift+S", "Save As"}, {"Mod+W", "Close document"},
	}},
	{"Edit", []shortcutHelp{
		{"Mod+Z", "Undo"}, {"Mod+Shift+Z", "Redo"}, {"Mod+C", "Copy"}, {"Mod+X", "Cut"}, {"Mod+V", "Paste"},
		{"Mod+D", "Duplicate selection or page"}, {"Delete", "Delete selection"}, {"Esc", "Cancel / deselect"},
		{"Up", "Move selection 1 pt"}, {"Shift+Up", "Move selection 10 pt"},
	}},
	{"View", []shortcutHelp{
		{"Mod+=", "Zoom in"}, {"Mod+-", "Zoom out"}, {"Mod+0", "Fit width"}, {"Mod+1", "Actual size"},
		{"Mod+B", "Sidebar"}, {"Mod+I", "Inspector"}, {"Mod+F", "Search"}, {"PageDown", "Next page"}, {"PageUp", "Previous page"},
	}},
	{"Pages", []shortcutHelp{
		{"Mod+R", "Rotate clockwise"}, {"Mod+Shift+R", "Rotate counter-clockwise"}, {"Mod+Backspace", "Delete page"},
	}},
	{"Tools", []shortcutHelp{
		{"V", "Select"}, {"T", "Text"}, {"P", "Pen"}, {"M", "Marker"}, {"H", "Highlight"},
		{"R", "Rectangle"}, {"E", "Ellipse"}, {"L", "Line"}, {"A", "Arrow"}, {"W", "Whiteout"}, {"S", "Signature"},
	}},
}
