package components

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"

	"gopdf/ui/theme"
)

// NewDialog returns a modal dialog in the app's style whose open state
// is *open.
func NewDialog(open *bool, content func() giode.Element) *dialog.Dialog {
	return giode.Dialog(open, content, giode.Styles{
		Background: theme.Surface, BorderRadius: theme.RadiusLarge,
		BoxShadow: giode.BoxShadow{Y: 24, Blur: 64, Color: "#0f172a40"},
		Padding:   giode.UniformInset(24),
	})
}
