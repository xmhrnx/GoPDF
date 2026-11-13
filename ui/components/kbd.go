package components

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// Kbd shows a keyboard shortcut spec ("Mod+S") the platform's way.
func Kbd(spec string) giode.Element {
	return giode.Box(giode.Styles{
		Padding: giode.SymmetricInset(6, 2), BorderRadius: 5, Background: theme.Subtle,
		BorderWidth: 1, BorderColor: theme.Border,
	}, giode.Text(giode.FormatShortcut(spec), giode.Styles{FontSize: 11, Color: theme.Muted, FontWeight: giode.FontWeightMedium}))
}
