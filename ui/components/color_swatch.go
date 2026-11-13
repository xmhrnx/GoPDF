package components

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// swatchSize is the diameter of a color swatch, in dp.
const swatchSize = 18

// ColorSwatch is a round color choice with a ring when selected. Create
// it once; call View every frame.
type ColorSwatch struct {
	color    string
	press    *giode.PressableElement
	selected func() bool
}

// NewColorSwatch returns a swatch of color (CSS hex) that calls onPick.
func NewColorSwatch(color string, selected func() bool, onPick func()) *ColorSwatch {
	s := &ColorSwatch{color: color, selected: selected}
	s.press = giode.Pressable(s.content).Label("Color " + color).OnClick(func(giode.ClickEvent) { onPick() })
	return s
}

// View returns the swatch for this frame.
func (s *ColorSwatch) View() giode.Element { return s.press }

func (s *ColorSwatch) content(st giode.PressState) giode.Element {
	ring := "transparent"
	switch {
	case s.selected():
		ring = theme.Primary
	case st.Hovered:
		ring = theme.Border
	}
	return giode.Box(giode.Styles{
		Width: swatchSize + 6, Height: swatchSize + 6, BorderRadius: 99,
		BorderWidth: 2, BorderColor: ring, Align: giode.AlignCenter, Justify: giode.JustifyCenter,
	}, giode.Box(giode.Styles{
		Width: swatchSize, Height: swatchSize, BorderRadius: 99, Background: s.color,
		BorderWidth: 1, BorderColor: theme.Border,
	}))
}
