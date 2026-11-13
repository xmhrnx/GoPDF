// Package assets embeds the app's icon, interface icons and fonts.
package assets

import "embed"

// AppIcon is the window and application icon (PNG).
//
//go:embed icon.png
var AppIcon []byte

// Icons holds the interface icons (Lucide, ISC license) under "icons".
//
//go:embed icons/*.svg
var Icons embed.FS

// Fonts holds the Inter typeface (SIL Open Font License) under "fonts".
//
//go:embed fonts/*.ttf
var Fonts embed.FS
