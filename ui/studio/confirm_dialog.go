package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"

	"gopdf/ui/components"
)

// confirmDialog asks before an action with consequences.
type confirmDialog struct {
	open    bool
	icon    string
	title   string
	message string
	then    func()
	ok      *giode.ButtonElement
	cancel  *giode.ButtonElement
	dialog  *dialog.Dialog
}

func newConfirmDialog() *confirmDialog {
	d := &confirmDialog{}
	d.ok = components.NewPrimaryButton("OK", func() {
		d.open = false
		d.then()
	})
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

// ask shows the question; then runs when action is chosen.
func (d *confirmDialog) ask(icon, title, message, action string, then func()) {
	d.icon, d.title, d.message, d.then, d.open = icon, title, message, then, true
	d.ok.Label(action)
}

func (d *confirmDialog) content() giode.Element {
	return components.DialogFrame(d.icon, d.title, d.message, giode.Box(giode.Styles{}), d.cancel, d.ok)
}
