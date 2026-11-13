package components

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// Section is a titled group of controls in a side panel.
func Section(title string, children ...giode.Element) giode.Element {
	kids := append([]giode.Element{giode.Text(title, theme.Overline)}, children...)
	return giode.Stack(giode.Styles{Gap: 10, Width: -1, Padding: giode.SymmetricInset(16, 14)}, kids...)
}
