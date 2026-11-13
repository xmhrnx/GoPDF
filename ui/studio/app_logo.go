package studio

import (
	"image"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"
)

// appLogo shows the app icon at size dp.
func appLogo(logo image.Image, size int) giode.Element {
	return giode.Image(logo, giode.Styles{Width: properties.Width(size), Height: properties.Height(size)})
}
