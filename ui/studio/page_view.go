package studio

import (
	"fmt"
	"math"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/internal/geom"
	"gopdf/internal/pagecache"
	"gopdf/ui/theme"
)

// maxRenderSide caps the longest side of a rendered page, in pixels, so
// deep zoom doesn't allocate hundreds of megabytes per page.
const maxRenderSide = 3200

// pageView is one page in the viewer: the rendered page, the markup
// drawn over it, and the pointer gestures of the active tool.
type pageView struct {
	s     *Studio
	page  int
	press *giode.PressableElement
}

func newPageView(s *Studio, page int) *pageView {
	pv := &pageView{s: s, page: page}
	pv.press = giode.Pressable(pv.content).Label(fmt.Sprintf("Page %d", page+1)).Focusable(false).
		OnClick(func(e giode.ClickEvent) { s.clickPage(page, pv.point(e.PositionDp.X, e.PositionDp.Y)) }).
		OnDragStart(func(e giode.DragEvent) { s.startGesture(page, pv.point(e.PositionDp.X, e.PositionDp.Y)) }).
		OnDragMove(func(e giode.DragEvent) { s.moveGesture(page, pv.point(e.PositionDp.X, e.PositionDp.Y)) }).
		OnDragEnd(func(e giode.DragEvent) { s.endGesture(e.Cancelled) })
	return pv
}

// point converts a position in dp on the page to page points.
func (pv *pageView) point(x, y float32) geom.Point {
	z := pv.s.effectiveZoom()
	return geom.Point{X: float64(x) / z, Y: float64(y) / z}
}

// View returns the page for this frame.
func (pv *pageView) View() giode.Element {
	pv.press.Styles(giode.Styles{Cursor: pv.s.cursor()})
	return pv.press
}

func (pv *pageView) content(giode.PressState) giode.Element {
	if pv.page >= pv.s.pageCount() {
		return giode.Box(giode.Styles{})
	}
	info := pv.s.doc.Doc().Pages[pv.page]
	z := pv.s.effectiveZoom()
	pw, ph := info.Size()
	w, h := int(math.Round(pw*z)), int(math.Round(ph*z))
	kids := []giode.Element{pv.image(w, h, pw, ph, z)}
	kids = append(kids, markupLayer(pv.s, pv.page, w, h, z))
	kids = append(kids, textMarkup(pv.s, pv.page, z)...)
	return giode.Box(giode.Styles{
		Width: properties.Width(w), Height: properties.Height(h), Background: "#ffffff",
		BoxShadow: giode.BoxShadow{Y: 2, Blur: 14, Color: theme.PageShadow},
	}, kids...)
}

func (pv *pageView) image(w, h int, pw, ph, z float64) giode.Element {
	dpi := 72 * z * float64(giode.DisplayScale())
	if longest := max(pw, ph) * dpi / 72; longest > maxRenderSide {
		dpi *= maxRenderSide / longest
	}
	key := pagecache.Key{Version: pv.s.doc.Version(), Page: pv.page, DPI: int(math.Round(dpi/4) * 4)}
	img, fallback, ok := pv.s.pages.Get(key)
	if !ok {
		img = fallback
	}
	if img == nil {
		return giode.Box(giode.Styles{Width: -1, Height: -1, Align: giode.AlignCenter, Justify: giode.JustifyCenter},
			giode.Spinner(giode.Styles{Width: 22, Color: "#94a3b8"}))
	}
	return giode.Image(img, giode.Styles{Width: properties.Width(w), Height: properties.Height(h), Fit: giode.FitStretch})
}
