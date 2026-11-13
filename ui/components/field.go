package components

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// Field stacks a label above a control.
func Field(label string, control giode.Element) giode.Element {
	return giode.Stack(giode.Styles{Gap: 6, Width: -1},
		giode.Text(label, theme.Label),
		control,
	)
}
