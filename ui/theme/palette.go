// Package theme holds the look of GoPDF: palettes, shared colors,
// sizes, text styles, fonts and icons.
package theme

import "github.com/xmhrnx/giode"

// Light is the light palette: cool neutrals with an indigo accent.
var Light = giode.Palette{
	Primary:    "#4f46e5",
	Secondary:  "#64748b",
	Background: "#f5f6f8",
	Surface:    "#ffffff",
	Text:       "#111827",
	Border:     "#e4e6eb",
	OnPrimary:  "#ffffff",
	Muted:      "#6b7280",
	Accent:     "#7c3aed",
	Success:    "#16a34a",
	Warning:    "#d97706",
	Danger:     "#e11d48",
	Info:       "#0284c7",
}

// Dark is the dark palette: near-black graphite with a softer indigo.
var Dark = giode.Palette{
	Primary:    "#818cf8",
	Secondary:  "#94a3b8",
	Background: "#0f1115",
	Surface:    "#171a20",
	Text:       "#e8eaef",
	Border:     "#262a33",
	OnPrimary:  "#0f1115",
	Muted:      "#8b93a3",
	Accent:     "#a78bfa",
	Success:    "#22c55e",
	Warning:    "#f59e0b",
	Danger:     "#fb7185",
	Info:       "#38bdf8",
}

// Apply installs both palettes and selects mode.
func Apply(mode giode.Mode) {
	giode.SetPalette(giode.ModeLight, Light)
	giode.SetPalette(giode.ModeDark, Dark)
	giode.SetMode(mode)
}

// ToggleMode switches between light and dark.
func ToggleMode() {
	if IsDark() {
		giode.SetMode(giode.ModeLight)
		return
	}
	giode.SetMode(giode.ModeDark)
}

// IsDark reports whether the dark palette is active.
func IsDark() bool { return giode.CurrentMode() == giode.ModeDark }
