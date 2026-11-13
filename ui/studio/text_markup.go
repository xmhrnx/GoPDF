package studio

import (
	"math"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/internal/annot"
)

// textMarkup lays out the page's text annotations at zoom z, each placed
// absolutely at its anchor.
func textMarkup(s *Studio, page int, z float64) []giode.Element {
	var out []giode.Element
	for _, a := range s.doc.Annots() {
		if a.Page != page || a.Kind != annot.Text {
			continue
		}
		size := a.FontSize * z
		out = append(out, giode.Box(giode.Styles{
			Position: giode.PositionAbsolute,
			Left:     properties.Left(edge(a.Points[0].X * z)), Top: properties.Top(edge(a.Points[0].Y * z)),
		}, giode.Text(a.Text, giode.Styles{
			FontSize: properties.FontSize(size), LineHeight: properties.LineHeight(a.LineHeight() * z),
			Color: a.Color.Hex(),
		})))
	}
	return out
}

// edge converts an offset in dp to an absolute position, keeping 0 set.
func edge(v float64) int {
	if r := int(math.Round(v)); r != 0 {
		return r
	}
	return giode.Flush
}
