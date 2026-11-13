package components

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// HRule is a full-width hairline separating panel sections.
func HRule() giode.Element {
	return giode.Box(giode.Styles{Width: -1, Height: 1, Background: theme.Border})
}
