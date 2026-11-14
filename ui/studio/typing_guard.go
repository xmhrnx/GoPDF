package studio

import "gopdf/internal/annot"

// editingText reports whether one of the app's text fields has the
// focus. Giode keeps typed characters, Delete and the arrows from
// reaching shortcuts while a field has focus, but not Mod shortcuts:
// Mod+C in a field must copy its text, not the selected markup.
func (v *views) editingText() bool {
	if a, ok := v.s.selectedAnnot(); ok && v.s.view.inspector {
		if a.Kind == annot.Text && v.inspector.text.field.Focused() {
			return true
		}
		if v.inspector.style.focused() {
			return true
		}
	}
	searchShown := v.s.view.sidebar && v.s.view.tab == tabSearch
	return searchShown && v.search.field.Focused()
}

func (v *views) dialogOpen() bool {
	return v.password.open || v.unsaved.open || v.properties.open || v.watermark.open ||
		v.pageNumbers.open || v.security.open || v.extract.open || v.split.open || v.shortcuts.open || v.signer.open
}

// unlessEditing wraps a clipboard shortcut so it leaves text fields and
// dialogs their own copy and paste.
func (v *views) unlessEditing(fn func()) func() {
	return func() {
		if !v.editingText() && !v.dialogOpen() {
			fn()
		}
	}
}

// unlessDialog wraps a shortcut that acts on the document behind a
// dialog, so it waits until the dialog is closed.
func (v *views) unlessDialog(fn func()) func() {
	return func() {
		if !v.dialogOpen() {
			fn()
		}
	}
}
