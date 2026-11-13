package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// dropOverlay invites a drop while files are dragged over the window.
func dropOverlay(visible bool) giode.Element {
	if !visible {
		return giode.Box(giode.Styles{Display: giode.DisplayNone})
	}
	return giode.Box(giode.Styles{
		Position: giode.PositionAbsolute, Top: giode.Flush, Left: giode.Flush, Right: giode.Flush, Bottom: giode.Flush,
		ZIndex: 50, Padding: giode.UniformInset(16), Background: giode.Mix("theme.background", "#000000", 0.25),
	}, giode.Box(giode.Styles{
		Width: -1, Height: -1, BorderRadius: theme.RadiusLarge, BorderWidth: 2, BorderColor: theme.Primary,
		Background: theme.Selected, Align: giode.AlignCenter, Justify: giode.JustifyCenter, FlexDirection: giode.FlexDirectionColumn, Gap: 12,
	},
		theme.Icon("upload", 40, theme.Primary),
		giode.Text("Drop to open", giode.Styles{FontSize: 20, Color: theme.Primary, FontWeight: giode.FontWeightSemiBold}),
	))
}
