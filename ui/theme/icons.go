package theme

import (
	"io/fs"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/icon"
	"github.com/xmhrnx/giode/styles/properties"
)

// RegisterIcons registers every SVG under dir in fsys by file name.
func RegisterIcons(fsys fs.FS, dir string) error {
	return giode.RegisterIconsFS(fsys, dir)
}

// Icon returns the named icon at size dp in color.
func Icon(name string, size int, color string) giode.Element {
	return giode.Icon(name, giode.Styles{Width: properties.Width(size), Height: properties.Height(size), Color: color})
}

// MenuIcon returns the registered icon name as a source for menu items;
// an unknown name gives an empty icon rather than a crash.
func MenuIcon(name string) icon.Source {
	if src, ok := icon.Lookup(name); ok {
		return src
	}
	return giode.IconPath{}
}
