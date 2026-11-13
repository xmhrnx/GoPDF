package studio

import (
	"github.com/xmhrnx/giode"

	"gopdf/ui/theme"
)

// featureCard advertises one capability on the welcome screen.
func featureCard(icon, title, text string) giode.Element {
	return giode.Stack(giode.Styles{
		Width: 146, Gap: 8, Padding: giode.UniformInset(14), BorderRadius: 12,
		Background: theme.Surface, BorderWidth: 1, BorderColor: theme.Border,
	},
		giode.Box(giode.Styles{Width: 30, Height: 30, BorderRadius: 8, Background: theme.Selected, Align: giode.AlignCenter, Justify: giode.JustifyCenter},
			theme.Icon(icon, 16, theme.Primary)),
		giode.Text(title, giode.Styles{FontSize: 13, Color: theme.TextColor, FontWeight: giode.FontWeightSemiBold}),
		giode.Text(text, theme.Caption),
	)
}
