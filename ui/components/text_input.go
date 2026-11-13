package components

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/ui/theme"
)

// TextInputStyles is the look of every text field.
var TextInputStyles = giode.Styles{
	Background: theme.Surface, Color: theme.TextColor, BorderColor: theme.Border,
	BorderRadius: theme.Radius, FontSize: 13, Padding: giode.SymmetricInset(10, 8),
}

// NewTextInput returns a single-line field in the app's style.
func NewTextInput(placeholder string) *input.Input {
	return giode.Input.Text(placeholder, TextInputStyles)
}
