package studio

import (
	"errors"
	"strings"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dropdown"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/internal/digitalid"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// The choices after the stored IDs in the "Signed by" menu.
const (
	signerNew  = "Someone new…"
	signerNone = "Nobody: just the mark, no digital ID"
)

// signerPanel is the "Signed by" part of the new signature dialog: who
// the signature belongs to, recorded in the documents it signs.
type signerPanel struct {
	s       *Studio
	choice  *dropdown.Select
	name    *input.Input
	email   *input.Input
	org     *input.Input
	pass    *input.Input
	confirm *input.Input
}

func newSignerPanel(s *Studio) *signerPanel {
	return &signerPanel{
		s:       s,
		choice:  giode.Input.Select(giode.Styles{Background: theme.Surface, Color: theme.TextColor, BorderColor: theme.Border, BorderRadius: theme.Radius, Width: -1}),
		name:    components.NewTextInput("Full name"),
		email:   components.NewTextInput("name@example.com").Keyboard(giode.KeyboardEmail),
		org:     components.NewTextInput("Organization (optional)"),
		pass:    components.NewTextInput("At least 6 characters").Password(true),
		confirm: components.NewTextInput("Repeat it").Password(true),
	}
}

// reset lists the stored IDs, with the one used last selected, or the
// form for someone new when there are none.
func (p *signerPanel) reset() {
	opts := make([]string, 0, len(p.s.ids.list)+2)
	selected := 0
	for i, info := range p.s.ids.list {
		label := info.Identity.Name
		if info.Identity.Email != "" {
			label += " · " + info.Identity.Email
		}
		opts = append(opts, label)
		if info.Fingerprint == p.s.defaultDigitalID() {
			selected = i
		}
	}
	p.choice.SetOptions(append(opts, signerNew, signerNone)...)
	p.choice.SetSelected(selected)
	for _, f := range []*input.Input{p.name, p.email, p.org, p.pass, p.confirm} {
		f.SetText("")
	}
}

func (p *signerPanel) isNew() bool  { return p.choice.Value() == signerNew }
func (p *signerPanel) isNone() bool { return p.choice.Value() == signerNone }

// signerName is the name of the person signing, for a typed signature.
func (p *signerPanel) signerName() string {
	if i := p.choice.Selected(); i < len(p.s.ids.list) {
		return p.s.ids.list[i].Identity.Name
	}
	if p.isNew() {
		return strings.TrimSpace(p.name.Text())
	}
	return ""
}

// resolve finds the digital ID chosen, creating it for someone new in the
// background, and hands its fingerprint ("" for none) to then, or the
// problem to fail.
func (p *signerPanel) resolve(then func(fingerprint string), fail func(error)) {
	switch i := p.choice.Selected(); {
	case i < len(p.s.ids.list):
		then(p.s.ids.list[i].Fingerprint)
		return
	case p.isNone():
		then("")
		return
	}
	if p.pass.Text() != p.confirm.Text() {
		fail(errors.New("the passwords don't match"))
		return
	}
	id := digitalid.Identity{Name: p.name.Text(), Email: p.email.Text(), Organization: p.org.Text()}
	pass := p.pass.Text()
	go func() {
		info, err := p.s.ids.store.Create(id, pass)
		p.s.win.Post(func() {
			if err != nil {
				fail(err)
				return
			}
			p.s.ids.added(info)
			then(info.Fingerprint)
		})
	}()
}

// View returns the panel for this frame.
func (p *signerPanel) View() giode.Element {
	kids := []giode.Element{components.Field("Signed by", p.choice)}
	switch {
	case p.isNew():
		kids = append(kids,
			giode.Row(giode.Styles{Width: -1, Gap: 8, Wrap: giode.WrapNoWrap},
				giode.Box(giode.Styles{FlexGrow: 1}, components.Field("Name", p.name)),
				giode.Box(giode.Styles{FlexGrow: 1}, components.Field("Email", p.email))),
			p.org,
			giode.Row(giode.Styles{Width: -1, Gap: 8, Wrap: giode.WrapNoWrap},
				giode.Box(giode.Styles{FlexGrow: 1}, components.Field("Password", p.pass)),
				giode.Box(giode.Styles{FlexGrow: 1}, components.Field("Confirm", p.confirm))),
			giode.Text("Creates your digital ID. The password protects it on this computer; you'll enter it when you save a signed document.", theme.Caption),
		)
	case p.isNone():
		kids = append(kids, giode.Text("The signature is only drawn on the page: readers can't check who signed.", theme.Caption))
	default:
		kids = append(kids, giode.Text("Your name, email and organization are sealed into every document you sign with this signature.", theme.Caption))
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 10}, kids...)
}
