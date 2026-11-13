package components

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/segmented"

	"gopdf/ui/theme"
)

// NewSegmented returns a full-width segmented control in the app's style.
func NewSegmented(options ...string) *segmented.Control {
	return giode.Input.Segmented(options...).Styles(giode.Styles{
		Width: -1, FontSize: 12, Background: theme.Subtle, BorderColor: theme.Surface, Color: theme.TextColor,
	})
}
