package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/segmented"

	"gopdf/internal/editor"
	"gopdf/internal/pdf"
	"gopdf/ui/components"
)

var (
	numberFormats   = []string{"%p", "Page %p", "Page %p of %P", "%p / %P"}
	numberPositions = []pdf.PageNumberPosition{pdf.BottomLeft, pdf.BottomCenter, pdf.BottomRight, pdf.TopLeft, pdf.TopCenter, pdf.TopRight}
)

// pageNumbersDialog stamps page numbers on every page.
type pageNumbersDialog struct {
	s        *Studio
	open     bool
	format   *segmented.Control
	vertical *segmented.Control
	align    *segmented.Control
	apply    *giode.ButtonElement
	cancel   *giode.ButtonElement
	dialog   *dialog.Dialog
}

func newPageNumbersDialog(s *Studio) *pageNumbersDialog {
	d := &pageNumbersDialog{s: s}
	d.format = components.NewSegmented("1", "Page 1", "Page 1 of 9", "1 / 9")
	d.format.SetSelected(2)
	d.vertical = components.NewSegmented("Bottom", "Top")
	d.align = components.NewSegmented("Left", "Center", "Right")
	d.align.SetSelected(1)
	d.apply = components.NewPrimaryButton("Add Numbers", d.submit)
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

func (d *pageNumbersDialog) submit() {
	d.open = false
	format := numberFormats[d.format.Selected()]
	pos := numberPositions[d.vertical.Selected()*3+d.align.Selected()]
	d.s.runEdit("Numbering pages", "Page numbers added", func(snap editor.Snapshot) (editor.Snapshot, error) {
		return editor.PageNumbers(snap, format, pos, 10)
	})
}

func (d *pageNumbersDialog) content() giode.Element {
	body := giode.Stack(giode.Styles{Width: -1, Gap: 14},
		components.Field("Format", d.format),
		components.Field("Position", d.vertical),
		components.Field("Alignment", d.align),
	)
	return components.DialogFrame("hash", "Add page numbers", "Numbers every page of the document.", body, d.cancel, d.apply)
}
