package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/segmented"

	"gopdf/internal/annot"
	"gopdf/internal/signature"
	"gopdf/ui/components"
)

// Signature dialog modes, in segmented-control order.
const (
	modeDraw = iota
	modeType
)

// signatureDialog creates a signature, drawn or typed, and saves it.
type signatureDialog struct {
	s      *Studio
	open   bool
	mode   *segmented.Control
	pad    *drawPad
	typed  *typePanel
	clear  *giode.ButtonElement
	cancel *giode.ButtonElement
	save   *giode.ButtonElement
	dialog *dialog.Dialog
}

func newSignatureDialog(s *Studio) *signatureDialog {
	d := &signatureDialog{s: s, pad: newDrawPad(), typed: newTypePanel(s)}
	d.mode = components.NewSegmented("Draw", "Type")
	d.clear = components.NewSecondaryButton("Clear", d.pad.clear)
	d.cancel = components.NewSecondaryButton("Cancel", func() { d.open = false })
	d.save = components.NewPrimaryButton("Save Signature", d.submit)
	d.dialog = components.NewDialog(&d.open, d.content)
	return d
}

func (d *signatureDialog) show() {
	d.pad.clear()
	d.typed.reset()
	d.open = true
	if d.mode.Selected() == modeType {
		d.typed.name.Focus()
	}
}

func (d *signatureDialog) submit() {
	var (
		sig signature.Signature
		ok  bool
	)
	if d.mode.Selected() == modeType {
		var err error
		if sig, err = d.typed.signature(); err != nil {
			d.s.toasts.Show(giode.ToastWarning, "Type your name first")
			return
		}
		ok = true
		d.s.memory.prefs.SignatureFont = d.typed.fontName()
	} else if sig, ok = d.pad.signature(); !ok {
		d.s.toasts.Show(giode.ToastWarning, "Draw your signature first")
		return
	}
	d.open = false
	d.s.signatures.add(d.s, sig)
	d.s.setTool(Tool(annot.Signature))
	d.s.toasts.Show(giode.ToastSuccess, "Signature saved. Click on a page to place it.")
}

func (d *signatureDialog) content() giode.Element {
	body := d.pad.View()
	actions := []giode.Element{d.clear, d.cancel, d.save}
	if d.mode.Selected() == modeType {
		body = d.typed.View()
		actions = actions[1:]
	}
	return components.DialogFrame("signature", "New signature",
		"Saved on this computer, ready to place in any document.",
		giode.Stack(giode.Styles{Width: -1, Gap: 14}, d.mode, body), actions...)
}
