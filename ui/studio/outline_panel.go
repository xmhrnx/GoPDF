package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// outlinePanel lists the document's bookmarks.
type outlinePanel struct {
	s    *Studio
	list *giode.LazyListElement
}

func newOutlinePanel(s *Studio) *outlinePanel {
	o := &outlinePanel{s: s}
	o.list = giode.LazyList(func() int { return len(s.outline) }, func(i int) any { return i }, o.row,
		giode.Styles{Gap: 1, ScrollBar: theme.ScrollBar})
	return o
}

// View returns the panel for this frame.
func (o *outlinePanel) View() giode.Element {
	if len(o.s.outline) == 0 {
		return components.EmptyHint("list-tree", "No outline", "This document has no bookmarks.")
	}
	return o.list
}

func (o *outlinePanel) row(i int) giode.Element {
	item := o.s.outline[i]
	p := giode.Keyed("outline", i, func() *giode.PressableElement {
		return giode.Pressable(func(st giode.PressState) giode.Element {
			it := o.s.outline[min(i, len(o.s.outline)-1)]
			bg := "transparent"
			if st.Hovered {
				bg = theme.Hover
			}
			weight := giode.FontWeightNormal
			if it.Level == 1 {
				weight = giode.FontWeightMedium
			}
			return giode.Row(giode.Styles{
				Width: -1, Gap: 8, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, BorderRadius: 6, Background: bg,
				Padding: giode.Inset{Top: 6, Bottom: 6, Right: 8, Left: 8 + 12*(it.Level-1)},
			},
				giode.Text(it.Title, giode.Styles{FontSize: 12.5, Color: theme.TextColor, FontWeight: weight, FlexGrow: 1, TextOverflow: giode.TextOverflowEllipsis}),
				giode.Text(fmt.Sprint(it.Page+1), theme.Caption),
			)
		}).OnClick(func(giode.ClickEvent) { o.s.goToPage(o.s.outline[i].Page) })
	})
	p.Label(item.Title)
	return giode.Box(giode.Styles{Width: -1, Padding: giode.SymmetricInset(8, 0)}, p)
}
