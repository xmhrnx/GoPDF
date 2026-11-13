package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// passwordDialog asks for the password of a protected file.
type passwordDialog struct {
	s       *Studio
	open    bool
	wrong   bool
	pending pendingOpen
	field   *input.Input
	ok      *giode.ButtonElement
	cancel  *giode.ButtonElement
	dialog  *dialog.Dialog
}

func newPasswordDialog(s *Studio) *passwordDialog {
	d := &passwordDialog{s: s}
	d.field = components.NewTextInput("Password").Password(true).Submit(func(string) { d.submit() })
	d.ok = components.NewPrimaryButton("Unlock", d.submit)
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

// ask shows the dialog for p; wrong says the last attempt failed.
func (d *passwordDialog) ask(p pendingOpen, wrong bool) {
	d.pending, d.wrong, d.open = p, wrong, true
	d.field.SetText("")
	d.field.Focus()
}

func (d *passwordDialog) submit() {
	d.open = false
	d.s.openBytes(d.pending.data, d.pending.path, d.pending.name, d.field.Text())
}

func (d *passwordDialog) content() giode.Element {
	body := []giode.Element{d.field}
	if d.wrong {
		body = append(body, giode.Text("That password isn't right. Try again.", giode.Styles{FontSize: 12, Color: theme.Danger}))
	}
	return components.DialogFrame("lock", "Password required",
		"“"+d.pending.name+"” is protected. Enter its password to open it.",
		giode.Stack(giode.Styles{Width: -1, Gap: 8}, body...), d.cancel, d.ok)
}
