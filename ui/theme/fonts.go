package theme

import (
	"io/fs"
	"path"

	"github.com/xmhrnx/giode"
)

// UIFont is the interface's typeface; main passes it as Config.Font, so
// every component uses it.
const UIFont = "Inter"

// uiFontFiles are the weights of UIFont in the font folder.
var uiFontFiles = []struct {
	file   string
	weight giode.FontWeight
}{
	{"Inter-Regular.ttf", giode.FontWeightNormal},
	{"Inter-Medium.ttf", giode.FontWeightMedium},
	{"Inter-SemiBold.ttf", giode.FontWeightSemiBold},
}

// RegisterFonts registers UIFont's weights from dir in fsys.
func RegisterFonts(fsys fs.FS, dir string) error {
	files := make([]giode.FontFile, 0, len(uiFontFiles))
	for _, f := range uiFontFiles {
		data, err := fs.ReadFile(fsys, path.Join(dir, f.file))
		if err != nil {
			return err
		}
		files = append(files, giode.FontFile{Family: UIFont, Weight: f.weight, Data: data})
	}
	return giode.RegisterFont(files...)
}
