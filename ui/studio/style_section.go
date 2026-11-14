package studio

import (
	"fmt"
	"strings"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/colorpicker"
	"github.com/xmhrnx/giode/components/numberinput"
	"github.com/xmhrnx/giode/components/segmented"
	"github.com/xmhrnx/giode/components/toggle"

	"gopdf/internal/annot"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

var strokeWidths = []float64{1, 2, 4, 8}

// Number field steps.
const (
	fontSizeStep      = 1
	signatureSizeStep = 10
)

// styleSection edits color, stroke, fill and sizes.
type styleSection struct {
	s          *Studio
	swatches   []*components.ColorSwatch
	custom     *components.HueSwatch
	picker     *colorpicker.ColorPicker
	showPicker bool
	width      *segmented.Control
	fill       *toggle.Switch
	font       *numberinput.NumberInput
	sigSize    *numberinput.NumberInput
	copy       *giode.ButtonElement
	remove     *giode.ButtonElement
}

func newStyleSection(s *Studio) *styleSection {
	st := &styleSection{s: s}
	for _, c := range annot.Swatches {
		st.swatches = append(st.swatches, components.NewColorSwatch(c.Hex(),
			func() bool { return st.color() == c }, func() { s.setColor(c) }))
	}
	presets := make([]string, len(annot.Swatches))
	for i, c := range annot.Swatches {
		presets[i] = c.Hex()
	}
	st.custom = components.NewHueSwatch(func() bool { return st.showPicker }, func() { st.showPicker = !st.showPicker })
	st.picker = giode.Input.Color(annot.DefaultStyle.Color.Hex(), giode.Styles{Width: -1}).Presets(presets...).
		OnChange(func(hex string) {
			if c, ok := annot.ParseHex(hex); ok {
				s.setColorLive(c)
			}
		})
	labels := make([]string, len(strokeWidths))
	for i, w := range strokeWidths {
		labels[i] = fmt.Sprintf("%g pt", w)
	}
	st.width = components.NewSegmented(labels...).OnChange(func(i int) { s.setStrokeWidth(strokeWidths[i]) })
	st.fill = giode.Input.Switch("Fill shape", giode.Styles{FontSize: 13, Color: theme.TextColor}).OnChange(s.setFilled)
	st.font = giode.Input.Number(annot.DefaultStyle.FontSize, giode.Styles{Width: -1, FontSize: 13}).
		Range(annot.MinFontSize, annot.MaxFontSize).Step(fontSizeStep).Decimals(1).Label("Font size").
		OnChange(s.setFontSize)
	st.sigSize = giode.Input.Number(annot.DefaultStyle.Size, giode.Styles{Width: -1, FontSize: 13}).
		Range(annot.MinSignatureSize, annot.MaxSignatureSize).Step(signatureSizeStep).Decimals(0).Label("Signature width").
		OnChange(s.setSignatureSize)
	st.copy = components.NewSecondaryButton("Duplicate", s.duplicateSelection).Icon(theme.Icon("copy", 15, ""))
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
		return annot.Style{Color: a.Color, Width: w, FontSize: a.FontSize, Filled: a.Filled, Size: a.Box().Dx()}
	}
	k, _ := st.s.styledKind()
	return st.s.memory.style(k)
}

func (st *styleSection) color() annot.Color { return st.current().Color }

// focused reports whether one of the section's number fields has focus.
func (st *styleSection) focused() bool {
	return st.font.Field().Focused() || st.sigSize.Field().Focused()
}

// blur takes the focus from the section's number fields.
func (st *styleSection) blur() {
	st.font.Field().Blur()
	st.sigSize.Field().Blur()
}

// colorRow is the swatches, the custom color wheel and, when open, the
// picker, kept in step with the current color.
func (st *styleSection) colorRow(cur annot.Color) giode.Element {
	views := make([]giode.Element, 0, len(st.swatches)+1)
	for _, sw := range st.swatches {
		views = append(views, sw.View())
	}
	views = append(views, st.custom.View())
	row := giode.Row(giode.Styles{Gap: 2, Width: -1}, views...)
	if !st.showPicker {
		return components.Field("Color", row)
	}
	if !strings.EqualFold(st.picker.Value(), cur.Hex()) {
		st.picker.SetValue(cur.Hex())
	}
	return components.Field("Color", giode.Stack(giode.Styles{Width: -1, Gap: 10}, row, st.picker))
}

// View returns the controls that apply to kind.
func (st *styleSection) View(kind annot.Kind) giode.Element {
	cur := st.current()
	var kids []giode.Element
	if kind != annot.Whiteout {
		kids = append(kids, st.colorRow(cur.Color))
	}
	if kind.Freehand() || kind.Linear() || kind == annot.Rect || kind == annot.Ellipse || kind == annot.Signature {
		st.width.SetSelected(nearest(strokeWidths, cur.Width))
		kids = append(kids, components.Field("Stroke", st.width))
	}
	if kind == annot.Rect || kind == annot.Ellipse {
		st.fill.SetChecked(cur.Filled)
		kids = append(kids, st.fill)
	}
	if kind == annot.Text {
		if !st.font.Field().Focused() {
			st.font.SetValue(cur.FontSize)
		}
		kids = append(kids, components.Field("Font size (pt)", st.font))
	}
	if kind == annot.Signature {
		if !st.sigSize.Field().Focused() {
			st.sigSize.SetValue(cur.Size)
		}
		kids = append(kids, components.Field("Width (pt)", st.sigSize))
	}
	if _, ok := st.s.selectedAnnot(); ok {
		kids = append(kids, giode.Row(giode.Styles{Gap: 8}, st.copy, st.remove))
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
