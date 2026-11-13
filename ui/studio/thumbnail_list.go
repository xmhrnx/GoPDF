package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// thumbnailList is the scrolling column of page thumbnails.
type thumbnailList struct {
	s    *Studio
	list *giode.LazyListElement
}

func newThumbnailList(s *Studio) *thumbnailList {
	t := &thumbnailList{s: s}
	t.list = giode.LazyList(s.pageCount, func(i int) any { return i }, t.row,
		giode.Styles{Gap: 4, ScrollBar: theme.ScrollBar})
	return t
}

// View returns the list for this frame.
func (t *thumbnailList) View() giode.Element { return t.list }

func (t *thumbnailList) row(i int) giode.Element {
	card := giode.Keyed("thumb", i, func() *thumbnailCard { return newThumbnailCard(t.s, i) })
	return giode.Row(giode.Styles{Width: -1, Justify: giode.JustifyCenter}, card.View())
}
