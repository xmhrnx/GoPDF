package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/ui/theme"
)

// logoBadge is the app mark: a document glyph on an indigo gradient.
func logoBadge(size int) giode.Element {
	return giode.Box(giode.Styles{
		Width: properties.Width(size), Height: properties.Height(size), BorderRadius: properties.BorderRadius(size * 3 / 10),
		BackgroundGradient: giode.Linear(135, "#6366f1", "#8b5cf6"),
		Align:              giode.AlignCenter, Justify: giode.JustifyCenter,
	}, theme.Icon("file-text", size*6/10, "#ffffff"))
}
