package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/internal/annot"
	"gopdf/internal/signature"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// Preview size of a saved signature, in dp.
const (
	sigCardWidth  = 200
	sigCardHeight = 56
)

// signatureCard is one saved signature: click to use it, × to delete it.
type signatureCard struct {
	s      *Studio
	sig    signature.Signature
	press  *giode.PressableElement
	remove *components.ToolButton
	active bool
}

func newSignatureCard(s *Studio, sig signature.Signature) *signatureCard {
	c := &signatureCard{s: s, sig: sig}
	c.press = giode.Pressable(c.content).Label("Use this signature").
		OnClick(func(giode.ClickEvent) {
			s.memory.setSignature(s, sig.ID)
			s.setTool(Tool(annot.Signature))
		})
	c.remove = components.NewToolButton("x", "Delete signature", func() { s.signatures.remove(s, sig.ID) }).Compact()
	return c
}

// View returns the card, highlighted when it is the one placed.
func (c *signatureCard) View(active bool) giode.Element {
	c.active = active
	return giode.Row(giode.Styles{Width: -1, Gap: 4, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		c.press, c.remove.View())
}

func (c *signatureCard) content(st giode.PressState) giode.Element {
	border := theme.Border
	switch {
	case c.active:
		border = theme.Primary
	case st.Hovered:
		border = theme.Faint
	}
	sig := c.sig
	return giode.Box(giode.Styles{
		Width: sigCardWidth, Height: sigCardHeight, Background: "#ffffff", BorderRadius: 8,
		BorderWidth: 1, BorderColor: border, Padding: giode.SymmetricInset(10, 6),
	}, giode.Canvas(giode.Styles{Width: -1, Height: -1}, func(p *giode.Painter) {
		w, h := p.Width(), p.Height()
		// Fit the signature's box into the card, keeping its shape.
		bw, bh := w, w/float32(sig.Aspect)
		if bh > h {
			bw, bh = h*float32(sig.Aspect), h
		}
		ox, oy := (w-bw)/2, (h-bh)/2
		for _, stroke := range sig.Strokes {
			if len(stroke) < 2 {
				continue
			}
			path := p.Path()
			path.MoveTo(ox+float32(stroke[0].X)*bw, oy+float32(stroke[0].Y)*bh)
			for _, u := range stroke[1:] {
				path.LineTo(ox+float32(u.X)*bw, oy+float32(u.Y)*bh)
			}
			path.Stroke(1.5, "#1e293b")
		}
	}))
}
