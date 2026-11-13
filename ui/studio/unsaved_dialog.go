package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"

	"gopdf/ui/components"
)

// unsavedDialog asks what to do with unsaved changes before they'd be lost.
type unsavedDialog struct {
	s       *Studio
	open    bool
	name    string
	then    func()
	save    *giode.ButtonElement
	discard *giode.ButtonElement
	cancel  *giode.ButtonElement
	dialog  *dialog.Dialog
}

func newUnsavedDialog(s *Studio) *unsavedDialog {
	d := &unsavedDialog{s: s}
	d.save = components.NewPrimaryButton("Save", func() {
		d.open = false
		s.saveThen(d.then)
	})
	d.discard = components.NewSecondaryButton("Don't Save", func() {
		d.open = false
		d.then()
	})
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

// ask shows the dialog for the document name; then runs once the
// changes are saved or discarded.
func (d *unsavedDialog) ask(name string, then func()) {
	d.name, d.then, d.open = name, then, true
}

func (d *unsavedDialog) content() giode.Element {
	return components.DialogFrame("save", "Save changes?",
		"“"+d.name+"” has changes that haven't been saved. They will be lost if you don't save them.",
		giode.Box(giode.Styles{}), d.discard, d.cancel, d.save)
}
