package studio

import (
	"gioui.org/io/key"
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// scrollSettleFrames is how long the current page ignores the scroll
// position after a jump, while the list applies it.
const scrollSettleFrames = 3

// viewer is the continuous, lazily built column of pages.
type viewer struct {
	s      *Studio
	list   *giode.LazyListElement
	settle int
	blur   bool // take the keyboard focus from text fields
}

func newViewer(s *Studio) *viewer {
	v := &viewer{s: s}
	v.list = giode.LazyList(s.pageCount, func(i int) any { return i }, v.row,
		giode.Styles{ScrollBar: theme.ScrollBar})
	return v
}

// scrollTo brings page to the top of the viewer.
func (v *viewer) scrollTo(page int) {
	v.list.ScrollToKey(page, giode.LazyAlignStart)
	v.settle = scrollSettleFrames
}

// View returns the viewer for this frame.
func (v *viewer) View() giode.Element {
	return giode.Box(giode.Styles{FlexGrow: 1, Height: -1, Background: theme.Canvas},
		giode.Raw(v.layout))
}

// blurFields takes the keyboard focus away from text fields, as a press
// on the page does in every editor.
func (v *viewer) blurFields() { v.blur = true }

func (v *viewer) layout(gtx giode.Context) giode.Dimensions {
	if v.blur {
		v.blur = false
		gtx.Execute(key.FocusCmd{})
	}
	v.s.view.viewportWidth = int(giode.PxToDp(gtx.Constraints.Max.X))
	dims := v.list.Layout(gtx)
	if v.settle > 0 {
		v.settle--
	} else if first, _ := v.list.VisibleRange(); first >= 0 && first < v.s.pageCount() {
		v.s.view.current = first
	}
	return dims
}

func (v *viewer) row(i int) giode.Element {
	pv := giode.Keyed("pageview", i, func() *pageView { return newPageView(v.s, i) })
	pad := giode.Inset{Top: theme.PageGap / 2, Bottom: theme.PageGap / 2}
	if i == 0 {
		pad.Top = theme.PageGap
	}
	if i == v.s.pageCount()-1 {
		pad.Bottom = theme.PageGap
	}
	return giode.Row(giode.Styles{Width: -1, Justify: giode.JustifyCenter, Wrap: giode.WrapNoWrap, Padding: pad}, pv.View())
}
