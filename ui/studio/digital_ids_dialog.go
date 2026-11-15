package studio

import (
	"context"
	"io"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"
	"github.com/xmhrnx/giode/components/segmented"

	"gopdf/internal/digitalid"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// Digital ID dialog modes.
const (
	idModeCreate = iota
	idModeImport
)

// digitalIDsDialog lists the digital IDs and creates or imports them.
type digitalIDsDialog struct {
	s        *Studio
	open     bool
	problem  string
	then     func(digitalid.Info) // runs with a new ID, when opened to get one
	mode     *segmented.Control
	name     *input.Input
	email    *input.Input
	org      *input.Input
	pass     *input.Input
	confirm  *input.Input
	file     []byte
	choose   *giode.ButtonElement
	fileName string
	done     *giode.ButtonElement
	add      *giode.ButtonElement
	dialog   *dialog.Dialog
}

func newDigitalIDsDialog(s *Studio) *digitalIDsDialog {
	d := &digitalIDsDialog{s: s}
	d.mode = components.NewSegmented("Create new", "Import .p12 / .pfx")
	d.name = components.NewTextInput("Full name, as it should appear")
	d.email = components.NewTextInput("name@example.com").Keyboard(giode.KeyboardEmail)
	d.org = components.NewTextInput("Optional")
	d.pass = components.NewTextInput("At least 6 characters").Password(true)
	d.confirm = components.NewTextInput("Repeat the password").Password(true)
	d.choose = components.NewSecondaryButton("Choose File…", d.pickFile).Icon(theme.Icon("folder-open", 15, ""))
	d.add = components.NewPrimaryButton("Add Digital ID", d.submit)
	d.done = components.NewSecondaryButton("Close", func() { d.open = false })
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

// show opens the manager; then, when not nil, receives the next ID added.
func (d *digitalIDsDialog) show(then func(digitalid.Info)) {
	d.then, d.problem, d.open = then, "", true
	d.mode.SetSelected(idModeCreate)
	for _, f := range []*input.Input{d.pass, d.confirm} {
		f.SetText("")
	}
	if len(d.s.ids.list) == 0 {
		d.name.Focus()
	}
}

func (d *digitalIDsDialog) pickFile() {
	go func() {
		files, err := d.s.win.PickFiles(context.Background(), giode.PickOptions{Types: []string{"p12", "pfx"}})
		if err != nil || len(files) == 0 {
			if err != nil {
				d.s.win.Post(func() { d.s.showError(err) })
			}
			return
		}
		data, err := io.ReadAll(files[0])
		files[0].Close()
		d.s.win.Post(func() {
			if err != nil {
				d.s.showError(err)
				return
			}
			d.file, d.fileName = data, files[0].Name
		})
	}()
}

func (d *digitalIDsDialog) submit() {
	if d.mode.Selected() == idModeCreate && d.pass.Text() != d.confirm.Text() {
		d.problem = "The passwords don't match."
		return
	}
	if d.mode.Selected() == idModeImport && d.file == nil {
		d.problem = "Choose the .p12 or .pfx file of your digital ID."
		return
	}
	id := digitalid.Identity{Name: d.name.Text(), Email: d.email.Text(), Organization: d.org.Text()}
	importing, file, pass := d.mode.Selected() == idModeImport, d.file, d.pass.Text()
	d.add.Loading(true)
	go func() {
		var (
			info digitalid.Info
			err  error
		)
		if importing {
			info, err = d.s.ids.store.Import(file, pass)
		} else {
			info, err = d.s.ids.store.Create(id, pass)
		}
		d.s.win.Post(func() {
			d.add.Loading(false)
			if err != nil {
				d.problem = err.Error()
				return
			}
			d.problem, d.file, d.fileName = "", nil, ""
			d.pass.SetText("")
			d.confirm.SetText("")
			d.s.ids.added(info)
			d.s.toasts.Show(giode.ToastSuccess, "Digital ID added for "+info.Identity.Name)
			if d.then != nil {
				d.open = false
				d.then(info)
			}
		})
	}()
}

func (d *digitalIDsDialog) content() giode.Element {
	var form giode.Element
	if d.mode.Selected() == idModeImport {
		file := "No file chosen"
		if d.fileName != "" {
			file = d.fileName
		}
		form = giode.Stack(giode.Styles{Width: -1, Gap: 12},
			giode.Row(giode.Styles{Gap: 10, Align: giode.AlignCenter}, d.choose, giode.Text(file, theme.Caption)),
			components.Field("Password of the file", d.pass),
			giode.Text("IDs from a certificate authority are trusted by other PDF readers.", theme.Caption),
		)
	} else {
		form = giode.Stack(giode.Styles{Width: -1, Gap: 12},
			components.Field("Name", d.name),
			components.Field("Email", d.email),
			components.Field("Organization", d.org),
			giode.Row(giode.Styles{Width: -1, Gap: 8, Wrap: giode.WrapNoWrap},
				giode.Box(giode.Styles{FlexGrow: 1}, components.Field("Password", d.pass)),
				giode.Box(giode.Styles{FlexGrow: 1}, components.Field("Confirm", d.confirm))),
			giode.Text("A self-signed ID: readers see who signed and that nothing changed, but can't verify your identity unless they trust your certificate.", theme.Caption),
		)
	}
	body := []giode.Element{idList(d.s), d.mode, form}
	if d.problem != "" {
		body = append(body, giode.Text(d.problem, giode.Styles{FontSize: 12, Color: theme.Danger}))
	}
	return components.DialogFrame("badge-check", "Digital IDs",
		"A digital ID signs documents in your name. It's stored on this computer, protected by its password.",
		giode.Stack(giode.Styles{Width: -1, Gap: 14}, body...), d.done, d.add)
}
