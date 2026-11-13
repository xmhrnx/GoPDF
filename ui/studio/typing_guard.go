package studio

import "gopdf/internal/annot"

// typing reports whether keys are meant as text: a field has the focus
// or a dialog is open. Shortcuts without modifiers (tool letters,
// Delete) stand down then, so typing never switches tools.
func (v *views) typing() bool {
	textShown := false
	if a, ok := v.s.selectedAnnot(); ok && a.Kind == annot.Text && v.s.view.inspector {
		textShown = true
	}
	searchShown := v.s.view.sidebar && v.s.view.tab == tabSearch
	return (textShown && v.inspector.text.field.Focused()) ||
		(searchShown && v.search.field.Focused()) || v.dialogOpen()
}

func (v *views) dialogOpen() bool {
	return v.password.open || v.unsaved.open || v.properties.open || v.watermark.open ||
		v.pageNumbers.open || v.security.open || v.extract.open || v.split.open || v.shortcuts.open || v.signaturePad.open
}

// unlessTyping wraps a plain-key shortcut action with the typing guard.
func (v *views) unlessTyping(fn func()) func() {
	return func() {
		if !v.typing() {
			fn()
		}
	}
}
