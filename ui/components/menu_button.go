package components

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"

	"gopdf/ui/theme"
)

// NewMenuButton returns a title-bar menu label ("File") that opens m.
func NewMenuButton(label string, m *menu.Menu) giode.Element {
	b := giode.Button(label, giode.Styles{
		Background: theme.Hover, Color: theme.TextColor, FontSize: 13,
		BorderRadius: 6, Padding: giode.SymmetricInset(10, 5),
	}).Variant(giode.VariantGhost)
	return giode.Row(giode.Styles{}, m.Trigger(b), m)
}
