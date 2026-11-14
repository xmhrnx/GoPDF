package components

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// hueSegments is how many slices draw the hue wheel.
const hueSegments = 24

// HueSwatch is a round hue wheel in the swatch row that opens a color
// picker. Create it once; call View every frame.
type HueSwatch struct {
	press  *giode.PressableElement
	active func() bool
}

// NewHueSwatch returns a wheel shown as pressed while active reports
// true, calling onClick when clicked.
func NewHueSwatch(active func() bool, onClick func()) *HueSwatch {
	h := &HueSwatch{active: active}
	h.press = giode.Pressable(h.content).Label("Custom color").OnClick(func(giode.ClickEvent) { onClick() })
	return h
}

// View returns the wheel for this frame.
func (h *HueSwatch) View() giode.Element { return h.press }

func (h *HueSwatch) content(st giode.PressState) giode.Element {
	ring := "transparent"
	switch {
	case h.active():
		ring = theme.Primary
	case st.Hovered:
		ring = theme.Border
	}
	return giode.Box(giode.Styles{
		Width: swatchSize + 6, Height: swatchSize + 6, BorderRadius: 99,
		BorderWidth: 2, BorderColor: ring, Align: giode.AlignCenter, Justify: giode.JustifyCenter,
	}, giode.Canvas(giode.Styles{Width: swatchSize, Height: swatchSize}, func(p *giode.Painter) {
		c := p.Width() / 2
		step := float32(360) / hueSegments
		for i := range hueSegments {
			p.Sector(c, c, 0, c, float32(i)*step, step+0.5, fmt.Sprintf("hsl(%d 85%% 55%%)", int(float32(i)*step)))
		}
		p.Circle(c, c, c*0.38, "#ffffff")
	}))
}
