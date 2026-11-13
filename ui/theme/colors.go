package theme

import "github.com/xmhrnx/giode"

// Theme tokens and opaque tints mixed from them, computed once: Mix with
// tokens follows light/dark switches.
var (
	Primary    = "theme.primary"
	OnPrimary  = "theme.onprimary"
	Surface    = "theme.surface"
	Background = "theme.background"
	TextColor  = "theme.text"
	Muted      = "theme.muted"
	Border     = "theme.border"
	Danger     = "theme.danger"
	Success    = "theme.success"

	// Canvas is the backdrop behind pages, a shade deeper than Background.
	Canvas = giode.Mix("theme.text", "theme.background", 0.06)
	// Hover tints a row or a ghost button under the pointer.
	Hover = giode.Mix("theme.text", "theme.surface", 0.06)
	// Pressed tints a control being pressed.
	Pressed = giode.Mix("theme.text", "theme.surface", 0.11)
	// Selected is the faint primary wash behind the active item.
	Selected = giode.Mix("theme.primary", "theme.surface", 0.12)
	// SelectedStrong is the primary wash of an active tool.
	SelectedStrong = giode.Mix("theme.primary", "theme.surface", 0.18)
	// Subtle is a field or chip background.
	Subtle = giode.Mix("theme.text", "theme.surface", 0.04)
	// Faint is text one step quieter than Muted.
	Faint = giode.Mix("theme.muted", "theme.surface", 0.35)
	// PageShadow is the drop shadow under pages.
	PageShadow = "#0f172a2e"
	// Selection outlines selected markup.
	Selection = "theme.primary"
)
