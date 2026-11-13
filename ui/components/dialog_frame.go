package components

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// dialogWidth is the width of every dialog panel, in dp.
const dialogWidth = 420

// DialogFrame lays out a dialog: an icon badge and title, a short
// explanation, the body and a row of actions aligned to the end.
func DialogFrame(icon, title, subtitle string, body giode.Element, actions ...giode.Element) giode.Element {
	header := giode.Row(giode.Styles{Gap: 12, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		giode.Box(giode.Styles{
			Width: 36, Height: 36, BorderRadius: 10, Background: theme.Selected,
			Align: giode.AlignCenter, Justify: giode.JustifyCenter,
		}, theme.Icon(icon, 18, theme.Primary)),
		giode.Stack(giode.Styles{Gap: 2, FlexGrow: 1},
			giode.Text(title, theme.Heading),
			giode.Text(subtitle, theme.BodyMute),
		),
	)
	footer := giode.Row(giode.Styles{Gap: 8, Justify: giode.JustifyEnd, Width: -1, Wrap: giode.WrapNoWrap}, actions...)
	return giode.Stack(giode.Styles{Width: dialogWidth, Gap: 18}, header, body, footer)
}
