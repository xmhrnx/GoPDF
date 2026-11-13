package annot

import "fmt"

// Color is an opaque RGB color.
type Color struct {
	R, G, B uint8
}

// Swatches are the colors offered in the inspector, in display order.
var Swatches = []Color{
	{0x11, 0x18, 0x27}, // ink
	{0xe1, 0x1d, 0x48}, // rose
	{0xf9, 0x73, 0x16}, // orange
	{0xfa, 0xcc, 0x15}, // yellow
	{0x22, 0xc5, 0x5e}, // green
	{0x06, 0xb6, 0xd4}, // cyan
	{0x3b, 0x82, 0xf6}, // blue
	{0x8b, 0x5c, 0xf6}, // violet
	{0xff, 0xff, 0xff}, // white
}

// White is the whiteout fill.
var White = Color{0xff, 0xff, 0xff}

// Hex returns the color as "#rrggbb".
func (c Color) Hex() string { return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B) }

// Floats returns the components in [0, 1], the form PDF operators take.
func (c Color) Floats() (r, g, b float64) {
	return float64(c.R) / 255, float64(c.G) / 255, float64(c.B) / 255
}
