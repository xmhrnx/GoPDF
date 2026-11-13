package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/internal/geom"
	"gopdf/internal/signature"
	"gopdf/ui/theme"
)

// Pad size and ink, in dp.
const (
	padWidth    = 420
	padHeight   = 170
	padInk      = 2.5
	padBaseline = 128
	padMargin   = 28
)

// drawPad is where a signature is drawn with the mouse, a pen or a finger.
type drawPad struct {
	strokes [][]geom.Point
	press   *giode.PressableElement
}

func newDrawPad() *drawPad {
	d := &drawPad{}
	d.press = giode.Pressable(d.content, giode.Styles{Cursor: giode.CursorCrosshair}).
		Label("Signature pad").Focusable(false).
		OnDragStart(func(e giode.DragEvent) {
			d.strokes = append(d.strokes, []geom.Point{{X: float64(e.PositionDp.X), Y: float64(e.PositionDp.Y)}})
		}).
		OnDragMove(func(e giode.DragEvent) {
			if n := len(d.strokes); n > 0 {
				d.strokes[n-1] = append(d.strokes[n-1], geom.Point{X: float64(e.PositionDp.X), Y: float64(e.PositionDp.Y)})
			}
		})
	return d
}

// clear wipes the pad.
func (d *drawPad) clear() { d.strokes = nil }

// signature returns what was drawn, or false when nothing was.
func (d *drawPad) signature() (signature.Signature, bool) {
	return signature.FromStrokes(d.strokes)
}

// View returns the pad for this frame.
func (d *drawPad) View() giode.Element {
	return giode.Stack(giode.Styles{Width: -1, Gap: 8}, d.press,
		giode.Text("Draw with the mouse, a pen or your finger.", theme.Caption))
}

func (d *drawPad) content(giode.PressState) giode.Element {
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
			p.Line(padMargin, padBaseline, padWidth-padMargin, padBaseline, 1, "#cbd5e1")
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
		giode.Box(giode.Styles{Position: giode.PositionAbsolute, Left: padMargin, Top: padBaseline + 8}, hint),
	)
}
