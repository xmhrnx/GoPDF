package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/internal/editor"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// minPasswordLen is the shortest password accepted.
const minPasswordLen = 4

// securityDialog sets the password the file is encrypted with on save.
type securityDialog struct {
	s       *Studio
	open    bool
	problem string
	user    *input.Input
	confirm *input.Input
	owner   *input.Input
	apply   *giode.ButtonElement
	remove  *giode.ButtonElement
	cancel  *giode.ButtonElement
	dialog  *dialog.Dialog
}

func newSecurityDialog(s *Studio) *securityDialog {
	d := &securityDialog{s: s}
	d.user = components.NewTextInput("Password to open the file").Password(true)
	d.confirm = components.NewTextInput("Repeat the password").Password(true)
	d.owner = components.NewTextInput("Optional, for changing permissions").Password(true)
	d.apply = components.NewPrimaryButton("Protect", d.submit)
	d.remove = components.NewSecondaryButton("Remove Password", func() {
		d.open = false
		d.setSecurity(editor.Security{})
	})
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

func (d *securityDialog) show() {
	sec := d.s.doc.Security
	d.user.SetText(sec.UserPassword)
	d.confirm.SetText(sec.UserPassword)
	d.owner.SetText(sec.OwnerPassword)
	d.problem, d.open = "", true
	d.user.Focus()
}

func (d *securityDialog) submit() {
	switch {
	case len(d.user.Text()) < minPasswordLen:
		d.problem = "Use at least 4 characters."
	case d.user.Text() != d.confirm.Text():
		d.problem = "The passwords don't match."
	default:
		d.open = false
		d.setSecurity(editor.Security{UserPassword: d.user.Text(), OwnerPassword: d.owner.Text()})
		d.s.toasts.Show(giode.ToastSuccess, "The file will be encrypted when you save it")
	}
}

func (d *securityDialog) setSecurity(sec editor.Security) {
	d.s.doc.Security = sec
	d.s.doc.Touch()
}

func (d *securityDialog) content() giode.Element {
	body := []giode.Element{
		components.Field("Password", d.user),
		components.Field("Confirm password", d.confirm),
		components.Field("Owner password", d.owner),
	}
	if d.problem != "" {
		body = append(body, giode.Text(d.problem, giode.Styles{FontSize: 12, Color: theme.Danger}))
	}
	actions := []giode.Element{d.cancel, d.apply}
	if d.s.hasDoc() && d.s.doc.Security.Enabled() {
		actions = append([]giode.Element{d.remove}, actions...)
	}
	return components.DialogFrame("lock", "Password protection",
		"Encrypts the file with AES-256 when you save. Anyone opening it will need the password.",
		giode.Stack(giode.Styles{Width: -1, Gap: 12}, body...), actions...)
}
