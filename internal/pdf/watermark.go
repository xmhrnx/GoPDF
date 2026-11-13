package pdf

import (
	"context"
	"fmt"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// Watermark describes a text stamp repeated on every page.
type Watermark struct {
	Text     string
	FontSize int
	Rotation int     // degrees, counter-clockwise
	Opacity  float64 // 0..1
	Color    string  // "#rrggbb"
	OnTop    bool
}

// PageNumberPosition is where page numbers go.
type PageNumberPosition string

// The supported positions, as pdfcpu anchors.
const (
	BottomCenter PageNumberPosition = "bc"
	BottomRight  PageNumberPosition = "br"
	BottomLeft   PageNumberPosition = "bl"
	TopCenter    PageNumberPosition = "tc"
	TopRight     PageNumberPosition = "tr"
	TopLeft      PageNumberPosition = "tl"
)

// AddWatermark stamps wm on every page.
func AddWatermark(data []byte, wm Watermark) ([]byte, error) {
	desc := fmt.Sprintf("font:Helvetica, points:%d, rot:%d, opacity:%.2f, fillcolor:%s, scale:1 abs",
		wm.FontSize, wm.Rotation, wm.Opacity, wm.Color)
	return stamp(data, wm.Text, desc, wm.OnTop)
}

// AddPageNumbers stamps "format" on every page, where %p is the page
// number and %P the page count ("Page %p of %P").
func AddPageNumbers(data []byte, format string, pos PageNumberPosition, fontSize int) ([]byte, error) {
	off := "0 24"
	if pos[0] == 't' {
		off = "0 -24"
	}
	desc := fmt.Sprintf("font:Helvetica, points:%d, rot:0, scale:1 abs, pos:%s, off:%s, fillcolor:#444444", fontSize, pos, off)
	return stamp(data, format, desc, true)
}

func stamp(data []byte, text, desc string, onTop bool) ([]byte, error) {
	return transform(data, func(c context.Context, rs io.ReadSeeker, w io.Writer, conf *model.Configuration) error {
		wm, err := api.TextWatermark(c, text, desc, onTop, false, types.POINTS, conf)
		if err != nil {
			return err
		}
		return api.AddWatermarks(c, rs, w, nil, wm, conf)
	})
}
