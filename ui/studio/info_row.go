package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// infoRow is a label and its value on one line.
func infoRow(label, value string) giode.Element {
	return giode.Row(giode.Styles{Width: -1, Gap: 12, Wrap: giode.WrapNoWrap, Justify: giode.JustifySpaceBetween},
		giode.Text(label, theme.BodyMute),
		giode.Text(value, giode.Styles{FontSize: 13, Color: theme.TextColor, MaxWidth: 160, TextOverflow: giode.TextOverflowEllipsis, TextAlign: giode.TextAlignEnd}),
	)
}
