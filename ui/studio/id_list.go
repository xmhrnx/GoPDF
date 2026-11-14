package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// idList shows the stored digital IDs, each with a delete button.
func idList(s *Studio) giode.Element {
	if len(s.ids.list) == 0 {
		return giode.Text("No digital IDs yet. Create one below, or import one you already have.", theme.Caption)
	}
	rows := make([]giode.Element, 0, len(s.ids.list))
	for _, info := range s.ids.list {
		rows = append(rows, giode.Keyed("digital-id", info.Fingerprint, func() *idRow { return newIDRow(s, info) }).View())
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 6}, rows...)
}
