package studio

import (
	"strings"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"
	"github.com/xmhrnx/giode/components/segmented"
	"github.com/xmhrnx/giode/components/toggle"

	"gopdf/internal/annot"
	"gopdf/internal/editor"
	"gopdf/internal/pdf"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

var (
	watermarkSizes     = []int{36, 60, 96}
	watermarkAngles    = []int{0, 45, 90}
	watermarkOpacities = []float64{0.15, 0.3, 0.5}
)

// watermarkDialog stamps text across every page.
type watermarkDialog struct {
	s        *Studio
	open     bool
	text     *input.Input
	size     *segmented.Control
	angle    *segmented.Control
	opacity  *segmented.Control
	onTop    *toggle.Switch
	color    annot.Color
	swatches []*components.ColorSwatch
	apply    *giode.ButtonElement
	cancel   *giode.ButtonElement
	dialog   *dialog.Dialog
}

func newWatermarkDialog(s *Studio) *watermarkDialog {
	d := &watermarkDialog{s: s, color: annot.Swatches[1]}
	d.text = components.NewTextInput("CONFIDENTIAL")
	d.size = components.NewSegmented("Small", "Medium", "Large")
	d.size.SetSelected(1)
	d.angle = components.NewSegmented("Flat", "Diagonal", "Vertical")
	d.angle.SetSelected(1)
	d.opacity = components.NewSegmented("Light", "Medium", "Strong")
	d.opacity.SetSelected(1)
	d.onTop = giode.Input.Switch("Draw over page content", giode.Styles{FontSize: 13, Color: theme.TextColor})
	d.onTop.SetChecked(true)
	for _, c := range annot.Swatches[:8] {
		d.swatches = append(d.swatches, components.NewColorSwatch(c.Hex(), func() bool { return d.color == c }, func() { d.color = c }))
	}
	d.apply = components.NewPrimaryButton("Add Watermark", d.submit)
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

func (d *watermarkDialog) show() {
	d.open = true
	d.text.Focus()
}

func (d *watermarkDialog) submit() {
	text := strings.TrimSpace(d.text.Text())
	if text == "" {
		text = "CONFIDENTIAL"
	}
	d.open = false
	wm := pdf.Watermark{
		Text: text, FontSize: watermarkSizes[d.size.Selected()], Rotation: watermarkAngles[d.angle.Selected()],
		Opacity: watermarkOpacities[d.opacity.Selected()], Color: d.color.Hex(), OnTop: d.onTop.Checked(),
	}
	d.s.runEdit("Adding watermark", "Watermark added", func(snap editor.Snapshot) (editor.Snapshot, error) {
		return editor.Watermark(snap, wm)
	})
}

func (d *watermarkDialog) content() giode.Element {
	swatches := make([]giode.Element, len(d.swatches))
	for i, sw := range d.swatches {
		swatches[i] = sw.View()
	}
	body := giode.Stack(giode.Styles{Width: -1, Gap: 14},
		components.Field("Text", d.text),
		components.Field("Size", d.size),
		components.Field("Angle", d.angle),
		components.Field("Opacity", d.opacity),
		components.Field("Color", giode.Row(giode.Styles{Gap: 2}, swatches...)),
		d.onTop,
	)
	return components.DialogFrame("stamp", "Add watermark", "Stamps the text on every page. Undo removes it.", body, d.cancel, d.apply)
}
