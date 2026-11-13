package components

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// buttonPadding is the padding of dialog and panel buttons.
var buttonPadding = giode.SymmetricInset(14, 8)

// NewPrimaryButton returns a filled call-to-action button.
func NewPrimaryButton(label string, onClick func()) *giode.ButtonElement {
	return giode.Button(label, giode.Styles{
		Background: theme.Primary, Color: theme.OnPrimary, BorderRadius: theme.Radius,
		FontSize: 13, Padding: buttonPadding,
	}).OnClick(onClick)
}

// NewSecondaryButton returns an outlined, quieter button.
func NewSecondaryButton(label string, onClick func()) *giode.ButtonElement {
	return giode.Button(label, giode.Styles{
		BorderColor: theme.Border, Background: theme.Hover, Color: theme.TextColor,
		BorderRadius: theme.Radius, FontSize: 13, Padding: buttonPadding,
	}).Variant(giode.VariantOutlined).OnClick(onClick)
}
