package pdf

import (
	"context"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"gopdf/internal/geom"
)

// PageInfo is the geometry of one page.
type PageInfo struct {
	// Box is the visible area (the crop box) in PDF user space.
	Box geom.Rect
	// Rotate is the page's clockwise display rotation: 0, 90, 180 or 270.
	Rotate int
}

// Size returns the page's display size in points, rotation applied.
func (p PageInfo) Size() (w, h float64) {
	w, h = p.Box.Dx(), p.Box.Dy()
	if p.Rotate%180 != 0 {
		w, h = h, w
	}
	return w, h
}

// displayMatrix returns the PDF matrix [a b c d e f] that maps display
// space (points, origin top-left, y down) to the page's user space.
func (p PageInfo) displayMatrix() [6]float64 {
	x0, y0, x1, y1 := p.Box.Min.X, p.Box.Min.Y, p.Box.Max.X, p.Box.Max.Y
	switch p.Rotate {
	case 90:
		return [6]float64{0, 1, 1, 0, x0, y0}
	case 180:
		return [6]float64{-1, 0, 0, 1, x1, y0}
	case 270:
		return [6]float64{0, -1, -1, 0, x1, y1}
	default:
		return [6]float64{1, 0, 0, -1, x0, y1}
	}
}

func normRotation(r int) int {
	r %= 360
	if r < 0 {
		r += 360
	}
	return r
}

// readPageInfos collects the geometry of every page in ctx.
func readPageInfos(ctx *model.Context) ([]PageInfo, error) {
	pages := make([]PageInfo, 0, ctx.PageCount)
	for n := 1; n <= ctx.PageCount; n++ {
		_, _, inh, err := ctx.PageDict(context.Background(), n, false)
		if err != nil {
			return nil, err
		}
		pages = append(pages, pageInfoFrom(inh))
	}
	return pages, nil
}

func pageInfoFrom(inh *model.InheritedPageAttrs) PageInfo {
	box := inh.CropBox
	if box == nil {
		box = inh.MediaBox
	}
	info := PageInfo{Rotate: normRotation(inh.Rotate)}
	if box == nil {
		info.Box = geom.Rect{Max: geom.Point{X: 612, Y: 792}} // US Letter
		return info
	}
	info.Box = geom.RectFromPoints(geom.Point{X: box.LL.X, Y: box.LL.Y}, geom.Point{X: box.UR.X, Y: box.UR.Y})
	return info
}
