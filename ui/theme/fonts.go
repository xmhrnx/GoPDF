package theme

import (
	"fmt"
	"io/fs"

	"gioui.org/font"
	"gioui.org/font/opentype"
	"github.com/xmhrnx/giode"
)

// RegisterFonts makes the Inter files in fsys (under dir) the only
// typeface, so every component, not only styled text, uses it.
func RegisterFonts(fsys fs.FS, dir string) error {
	files, err := fs.Glob(fsys, dir+"/*.ttf")
	if err != nil {
		return err
	}
	var faces []font.FontFace
	for _, name := range files {
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return err
		}
		parsed, err := opentype.ParseCollection(data)
		if err != nil {
			return fmt.Errorf("font %s: %w", name, err)
		}
		faces = append(faces, parsed...)
	}
	giode.RegisterFonts(faces)
	return nil
}
