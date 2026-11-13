package components

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/ui/theme"
)

// VDivider is a thin vertical rule height dp tall, for toolbars.
func VDivider(height int) giode.Element {
	return giode.Box(giode.Styles{
		Width: 1, Height: properties.Height(height), Background: theme.Border,
		Margin: giode.SymmetricInset(6, 0),
	})
}
