package components

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// VRule is a full-height hairline between panels.
func VRule() giode.Element {
	return giode.Box(giode.Styles{Width: 1, Height: -1, Background: theme.Border})
}
