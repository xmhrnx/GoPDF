package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/internal/digitalid"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// unlockIDDialog asks for a digital ID's password before signing.
type unlockIDDialog struct {
	s       *Studio
	open    bool
	info    digitalid.Info
	then    func()
	problem string
	field   *input.Input
	ok      *giode.ButtonElement
	cancel  *giode.ButtonElement
	dialog  *dialog.Dialog
}

func newUnlockIDDialog(s *Studio) *unlockIDDialog {
	d := &unlockIDDialog{s: s}
	d.field = components.NewTextInput("Password").Password(true).Submit(func(string) { d.submit() })
	d.ok = components.NewPrimaryButton("Unlock and Sign", d.submit)
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

// ask shows the dialog for info; then runs once the ID is unlocked.
func (d *unlockIDDialog) ask(info digitalid.Info, then func()) {
	d.info, d.then, d.problem, d.open = info, then, "", true
	d.field.SetText("")
	d.field.Focus()
}

func (d *unlockIDDialog) submit() {
	fp, pw := d.info.Fingerprint, d.field.Text()
	d.ok.Loading(true)
	go func() {
		signer, err := d.s.ids.store.Unlock(fp, pw)
		d.s.win.Post(func() {
			d.ok.Loading(false)
			if err != nil {
				d.problem = err.Error()
				d.field.SelectAll()
				return
			}
			d.s.ids.unlocked[fp] = signer
			d.open = false
			d.then()
		})
	}()
}

func (d *unlockIDDialog) content() giode.Element {
	body := []giode.Element{d.field}
	if d.problem != "" {
		body = append(body, giode.Text(d.problem, giode.Styles{FontSize: 12, Color: theme.Danger}))
	}
	return components.DialogFrame("key-round", "Sign as "+d.info.Identity.Name,
		"Enter the password of this digital ID. It stays unlocked until you quit GoPDF.",
		giode.Stack(giode.Styles{Width: -1, Gap: 8}, body...), d.cancel, d.ok)
}
