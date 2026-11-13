package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"

	"gopdf/internal/annot"
	"gopdf/internal/geom"
	"gopdf/internal/signature"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// Pad size and ink, in dp.
const (
	padWidth    = 420
	padHeight   = 170
	padInk      = 2.5
	padBaseline = 128
)

// signaturePadDialog lets the user draw a signature and save it.
type signaturePadDialog struct {
	s       *Studio
	open    bool
	strokes [][]geom.Point
	pad     *giode.PressableElement
	clear   *giode.ButtonElement
	cancel  *giode.ButtonElement
	save    *giode.ButtonElement
	dialog  *dialog.Dialog
}

func newSignaturePadDialog(s *Studio) *signaturePadDialog {
	d := &signaturePadDialog{s: s}
	d.pad = giode.Pressable(d.padContent, giode.Styles{Cursor: giode.CursorCrosshair}).
		Label("Signature pad").Focusable(false).
		OnDragStart(func(e giode.DragEvent) {
			d.strokes = append(d.strokes, []geom.Point{{X: float64(e.PositionDp.X), Y: float64(e.PositionDp.Y)}})
		}).
		OnDragMove(func(e giode.DragEvent) {
			if n := len(d.strokes); n > 0 {
				d.strokes[n-1] = append(d.strokes[n-1], geom.Point{X: float64(e.PositionDp.X), Y: float64(e.PositionDp.Y)})
			}
		})
	d.clear = components.NewSecondaryButton("Clear", func() { d.strokes = nil })
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.save = components.NewPrimaryButton("Save Signature", d.submit)
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

func (d *signaturePadDialog) show() {
	d.strokes, d.open = nil, true
}

func (d *signaturePadDialog) submit() {
	sig, ok := signature.FromStrokes(d.strokes)
	if !ok {
		d.s.toasts.Show(giode.ToastWarning, "Draw your signature first")
		return
	}
	d.open = false
	d.s.signatures.add(d.s, sig)
	d.s.setTool(Tool(annot.Signature))
	d.s.toasts.Show(giode.ToastSuccess, "Signature saved. Click on a page to place it.")
}

func (d *signaturePadDialog) content() giode.Element {
	body := giode.Stack(giode.Styles{Width: -1, Gap: 8}, d.pad,
		giode.Text("Draw with the mouse, a pen or your finger.", theme.Caption))
	return components.DialogFrame("signature", "New signature",
		"Saved on this computer, ready to place in any document.", body, d.clear, d.cancel, d.save)
}

func (d *signaturePadDialog) padContent(giode.PressState) giode.Element {
	strokes := d.strokes
	var hint giode.Element = giode.Box(giode.Styles{})
	if len(strokes) == 0 {
		hint = giode.Text("Sign here", giode.Styles{FontSize: 13, Color: "#94a3b8"})
	}
	return giode.Box(giode.Styles{
		Width: padWidth, Height: padHeight, Background: "#ffffff", BorderRadius: 10,
		BorderWidth: 1, BorderColor: theme.Border, Overflow: giode.OverflowHidden,
	},
		giode.Canvas(giode.Styles{Width: padWidth, Height: padHeight}, func(p *giode.Painter) {
			p.Line(28, padBaseline, padWidth-28, padBaseline, 1, "#cbd5e1")
			for _, st := range strokes {
				if len(st) < 2 {
					continue
				}
				path := p.Path()
				path.MoveTo(float32(st[0].X), float32(st[0].Y))
				for _, pt := range st[1:] {
					path.LineTo(float32(pt.X), float32(pt.Y))
				}
				path.Stroke(padInk, "#1e293b")
			}
		}),
		giode.Box(giode.Styles{Position: giode.PositionAbsolute, Left: 28, Top: padBaseline + 8}, hint),
	)
}
