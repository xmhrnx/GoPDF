package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/internal/editor"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// extractDialog saves chosen pages as a new PDF.
type extractDialog struct {
	s       *Studio
	open    bool
	problem string
	ranges  *input.Input
	export  *giode.ButtonElement
	cancel  *giode.ButtonElement
	dialog  *dialog.Dialog
}

func newExtractDialog(s *Studio) *extractDialog {
	d := &extractDialog{s: s}
	d.ranges = components.NewTextInput("e.g. 1-3, 5, 8-").Submit(func(string) { d.submit() })
	d.export = components.NewPrimaryButton("Extract", d.submit)
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

func (d *extractDialog) show(pages string) {
	d.ranges.SetText(pages)
	d.problem, d.open = "", true
	d.ranges.Focus()
}

func (d *extractDialog) submit() {
	pages, err := editor.ParseRanges(d.ranges.Text(), d.s.pageCount())
	if err != nil {
		d.problem = err.Error()
		return
	}
	d.open = false
	snap := d.s.doc.Snapshot()
	name := fmt.Sprintf("%s-pages.pdf", baseName(d.s.doc.Name))
	d.s.exportFile("Extracting pages", name, func() ([]byte, error) { return editor.ExportPages(snap, pages) })
}

func (d *extractDialog) content() giode.Element {
	body := []giode.Element{components.Field(fmt.Sprintf("Pages (1–%d)", d.s.pageCount()), d.ranges)}
	if d.problem != "" {
		body = append(body, giode.Text(d.problem, giode.Styles{FontSize: 12, Color: theme.Danger}))
	}
	return components.DialogFrame("file-output", "Extract pages",
		"Saves the chosen pages, with their markup, as a new PDF. Pages are kept in the order you list them.",
		giode.Stack(giode.Styles{Width: -1, Gap: 8}, body...), d.cancel, d.export)
}
