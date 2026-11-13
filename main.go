// Command gopdf is a fast, native PDF editor built with Giode.
package main

import (
	"bytes"
	"image"
	"image/png"
	"log"
	"os"
	"runtime/debug"

	"github.com/xmhrnx/giode"

	"gopdf/assets"
	"gopdf/ui/studio"
	"gopdf/ui/theme"
)

// gcPercent is the garbage collector's target heap growth.
const gcPercent = 50

func main() {
	// Rendered pages are large and short-lived: collect garbage sooner
	// than Go's default so the app stays light on memory.
	debug.SetGCPercent(gcPercent)
	if err := theme.RegisterFonts(assets.Fonts, "fonts"); err != nil {
		log.Fatal(err)
	}
	if err := theme.RegisterIcons(assets.Icons, "icons"); err != nil {
		log.Fatal(err)
	}
	scripts, err := theme.RegisterScriptFonts(assets.SignatureFonts, "signature-fonts")
	if err != nil {
		log.Fatal(err)
	}
	theme.Apply(giode.ModeLight)

	logo, err := png.Decode(bytes.NewReader(assets.AppIcon))
	if err != nil {
		log.Fatal(err)
	}
	app := studio.New(logo, scripts)
	win := giode.New(giode.Config{
		Title:      "GoPDF",
		Size:       image.Pt(1360, 860),
		MinSize:    image.Pt(900, 560),
		Resizable:  true,
		Frameless:  true,
		Center:     true,
		Background: "theme.background",
		Icon:       assets.AppIcon,
		OnFileDrop: app.OnFileDrop,
		OnFileDrag: app.OnFileDrag,
		OnOpenURL:  app.OnOpenURL,
		OnReady:    func() { app.OpenFiles(os.Args[1:]) },
	})
	app.Attach(win)
	if err := win.Run(app.View); err != nil {
		log.Fatal(err)
	}
}
