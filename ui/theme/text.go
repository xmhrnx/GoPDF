package theme

import "github.com/xmhrnx/giode"

// Text styles.
var (
	Body     = giode.Styles{FontSize: 13, Color: TextColor}
	BodyMute = giode.Styles{FontSize: 13, Color: Muted}
	Caption  = giode.Styles{FontSize: 11.5, Color: Muted}
	Label    = giode.Styles{FontSize: 12, Color: TextColor, FontWeight: giode.FontWeightMedium}
	Overline = giode.Styles{FontSize: 10.5, Color: Muted, FontWeight: giode.FontWeightSemiBold, TextTransform: giode.TextTransformUppercase}
	Heading  = giode.Styles{FontSize: 16, Color: TextColor, FontWeight: giode.FontWeightSemiBold}
	Display  = giode.Styles{FontSize: 30, Color: TextColor, FontWeight: giode.FontWeightSemiBold}
)
