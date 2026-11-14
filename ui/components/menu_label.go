package components

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// MenuLabel returns a title-bar menu label ("File"), for a menu bar's
// Add.
func MenuLabel(label string) giode.Element {
	return giode.Button(label, giode.Styles{
		Background: theme.Hover, Color: theme.TextColor, FontSize: 13,
		BorderRadius: 6, Padding: giode.SymmetricInset(10, 5),
	}).Variant(giode.VariantGhost)
}
