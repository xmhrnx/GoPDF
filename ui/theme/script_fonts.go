package theme

import (
	"io/fs"
	"path"

	"github.com/xmhrnx/giode"
)

// ScriptFont is a handwriting typeface offered for typed signatures.
type ScriptFont struct {
	Name string // display name, also the family text is styled with
	Data []byte // the TrueType file
}

// scriptFontFiles names the files in the signature font folder, in the
// order they are offered.
var scriptFontFiles = []struct{ name, file string }{
	{"Great Vibes", "GreatVibes-Regular.ttf"},
	{"Allura", "Allura-Regular.ttf"},
	{"Alex Brush", "AlexBrush-Regular.ttf"},
	{"Sacramento", "Sacramento-Regular.ttf"},
	{"Mrs Saint Delafield", "MrsSaintDelafield-Regular.ttf"},
	{"Cedarville Cursive", "Cedarville-Cursive.ttf"},
}

// RegisterScriptFonts loads the signature typefaces from dir in fsys and
// registers each under its name, so text can preview in it.
func RegisterScriptFonts(fsys fs.FS, dir string) ([]ScriptFont, error) {
	fonts := make([]ScriptFont, 0, len(scriptFontFiles))
	files := make([]giode.FontFile, 0, len(scriptFontFiles))
	for _, f := range scriptFontFiles {
		data, err := fs.ReadFile(fsys, path.Join(dir, f.file))
		if err != nil {
			return nil, err
		}
		fonts = append(fonts, ScriptFont{Name: f.name, Data: data})
		files = append(files, giode.FontFile{Family: f.name, Data: data})
	}
	return fonts, giode.RegisterFont(files...)
}
