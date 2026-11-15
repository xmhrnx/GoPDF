package studio

import (
	"errors"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/segmented"

	"gopdf/internal/annot"
	"gopdf/internal/signature"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// Signature dialog modes, in segmented-control order.
const (
	modeDraw = iota
	modeType
)

// signatureDialog creates a signature, drawn or typed, and saves it.
type signatureDialog struct {
	s       *Studio
	open    bool
	then    func() // runs once a signature is saved; nil to just pick the tool
	signer  *signerPanel
	problem string
	mode    *segmented.Control
	pad     *drawPad
	typed   *typePanel
	clear   *giode.ButtonElement
	cancel  *giode.ButtonElement
	save    *giode.ButtonElement
	dialog  *dialog.Dialog
}

func newSignatureDialog(s *Studio) *signatureDialog {
	d := &signatureDialog{s: s, signer: newSignerPanel(s), pad: newDrawPad(), typed: newTypePanel(s)}
	d.mode = components.NewSegmented("Draw", "Type")
	d.clear = components.NewSecondaryButton("Clear", d.pad.clear)
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.save = components.NewPrimaryButton("Save Signature", d.submit)
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

// show opens the dialog; then, when not nil, runs once the new signature
// is saved (to place it where the user clicked).
func (d *signatureDialog) show(then func()) {
	d.then, d.problem = then, ""
	d.signer.reset()
	d.pad.clear()
	d.typed.reset()
	d.open = true
	if d.mode.Selected() == modeType {
		d.typed.name.Focus()
	}
}

func (d *signatureDialog) submit() {
	sig, err := d.mark()
	if err != nil {
		d.problem = err.Error()
		return
	}
	d.save.Loading(true)
	d.signer.resolve(func(fingerprint string) {
		d.save.Loading(false)
		sig.DigitalID = fingerprint
		if fingerprint != "" {
			d.s.memory.prefs.DigitalID = fingerprint
		}
		d.finish(sig)
	}, func(err error) {
		d.save.Loading(false)
		d.problem = err.Error()
	})
}

// mark returns the drawn or typed signature.
func (d *signatureDialog) mark() (signature.Signature, error) {
	if d.mode.Selected() == modeType {
		sig, err := d.typed.signature()
		if err != nil {
			return sig, errors.New("type the name to sign with")
		}
		d.s.memory.prefs.SignatureFont = d.typed.fontName()
		return sig, nil
	}
	sig, ok := d.pad.signature()
	if !ok {
		return sig, errors.New("draw your signature in the box")
	}
	return sig, nil
}

// finish saves sig and places it, or picks the tool to place it with.
func (d *signatureDialog) finish(sig signature.Signature) {
	d.open = false
	d.s.signatures.add(d.s, sig)
	if d.then != nil {
		d.then()
		return
	}
	d.s.setTool(Tool(annot.Signature))
	d.s.toasts.Show(giode.ToastSuccess, "Signature saved. Click on a page to place it.")
}

func (d *signatureDialog) content() giode.Element {
	body := d.pad.View()
	actions := []giode.Element{d.clear, d.cancel, d.save}
	if d.mode.Selected() == modeType {
		if d.typed.name.Text() == "" && !d.typed.name.Focused() {
			d.typed.name.SetText(d.signer.signerName())
		}
		body = d.typed.View()
		actions = actions[1:]
	}
	kids := []giode.Element{d.signer.View(), components.HRule(), d.mode, body}
	if d.problem != "" {
		kids = append(kids, giode.Text(d.problem, giode.Styles{FontSize: 12, Color: theme.Danger}))
	}
	return components.DialogFrame("signature", "New signature",
		"Who is signing, and how the signature looks. Saved on this computer, ready for any document.",
		giode.Stack(giode.Styles{Width: -1, Gap: 14}, kids...), actions...)
}
