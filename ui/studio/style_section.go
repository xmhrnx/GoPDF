package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dropdown"
	"github.com/xmhrnx/giode/components/segmented"
	"github.com/xmhrnx/giode/components/toggle"

	"gopdf/internal/annot"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

var (
	strokeWidths = []float64{1, 2, 4, 8}
	fontSizes    = []float64{10, 12, 14, 16, 20, 24, 32, 48, 64}
)

// styleSection edits color, stroke, fill and font size.
type styleSection struct {
	s        *Studio
	swatches []*components.ColorSwatch
	width    *segmented.Control
	fill     *toggle.Switch
	font     *dropdown.Select
	remove   *giode.ButtonElement
}

func newStyleSection(s *Studio) *styleSection {
	st := &styleSection{s: s}
	for _, c := range annot.Swatches {
		st.swatches = append(st.swatches, components.NewColorSwatch(c.Hex(),
			func() bool { return st.color() == c }, func() { s.setColor(c) }))
	}
	labels := make([]string, len(strokeWidths))
	for i, w := range strokeWidths {
		labels[i] = fmt.Sprintf("%g pt", w)
	}
	st.width = components.NewSegmented(labels...).OnChange(func(i int) { s.setStrokeWidth(strokeWidths[i]) })
	st.fill = giode.Input.Switch("Fill shape", giode.Styles{FontSize: 13, Color: theme.TextColor}).OnChange(s.setFilled)
	sizes := make([]string, len(fontSizes))
	for i, f := range fontSizes {
		sizes[i] = fmt.Sprintf("%g pt", f)
	}
	st.font = giode.Input.Select(giode.Styles{Background: theme.Surface, Color: theme.TextColor, BorderColor: theme.Border, BorderRadius: theme.Radius, Width: -1}, sizes...).
		OnChange(func(i int) { s.setFontSize(fontSizes[i]) })
	st.remove = components.NewSecondaryButton("Delete", s.deleteSelected).Icon(theme.Icon("trash-2", 15, ""))
	return st
}

// current returns the style values shown: the selection's or the tool's.
func (st *styleSection) current() annot.Style {
	if a, ok := st.s.selectedAnnot(); ok {
		w := a.Width
		if a.Kind == annot.Marker {
			w = a.Width / 5
		}
		return annot.Style{Color: a.Color, Width: w, FontSize: a.FontSize, Filled: a.Filled}
	}
	return st.s.toolStyle(st.s.tool)
}

func (st *styleSection) color() annot.Color { return st.current().Color }

// View returns the controls that apply to kind.
func (st *styleSection) View(kind annot.Kind) giode.Element {
	cur := st.current()
	var kids []giode.Element
	if kind != annot.Whiteout {
		views := make([]giode.Element, len(st.swatches))
		for i, sw := range st.swatches {
			views[i] = sw.View()
		}
		kids = append(kids, components.Field("Color", giode.Row(giode.Styles{Gap: 2, Width: -1}, views...)))
	}
	if kind.Freehand() || kind.Linear() || kind == annot.Rect || kind == annot.Ellipse {
		st.width.SetSelected(nearest(strokeWidths, cur.Width))
		kids = append(kids, components.Field("Stroke", st.width))
	}
	if kind == annot.Rect || kind == annot.Ellipse {
		st.fill.SetChecked(cur.Filled)
		kids = append(kids, st.fill)
	}
	if kind == annot.Text {
		st.font.SetSelected(nearest(fontSizes, cur.FontSize))
		kids = append(kids, components.Field("Font size", st.font))
	}
	if _, ok := st.s.selectedAnnot(); ok {
		kids = append(kids, st.remove)
	}
	if len(kids) == 0 {
		kids = append(kids, giode.Text("Covers content with white when the file is saved. It hides the content from view; it does not remove it.", theme.Caption))
	}
	return components.Section("Style", kids...)
}

// nearest returns the index of the value in vs closest to v.
func nearest(vs []float64, v float64) int {
	best := 0
	for i := range vs {
		if abs(vs[i]-v) < abs(vs[best]-v) {
			best = i
		}
	}
	return best
}

func abs(v float64) float64 { return max(v, -v) }
