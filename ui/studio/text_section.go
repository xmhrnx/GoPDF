package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/internal/annot"
	"gopdf/ui/components"
)

// textSection edits the content of the selected text annotation.
type textSection struct {
	s       *Studio
	field   *input.Input
	fresh   bool // focused, nothing typed yet: the first change opens an undo step
	focusIn int  // frames until the field, new in the layout, takes the focus
}

func newTextSection(s *Studio) *textSection {
	t := &textSection{s: s}
	t.field = giode.Input.Text("Type here…", components.TextInputStyles).Lines(3, 10).
		SubmitOn(giode.SubmitNever).
		OnFocus(func() { t.fresh = true }).
		OnChange(t.changed)
	return t
}

func (t *textSection) changed(text string) {
	a, ok := t.s.selectedAnnot()
	if !ok || a.Kind != annot.Text || !t.s.idle() {
		return
	}
	if t.fresh {
		t.s.doc.BeginChange()
		t.fresh = false
	}
	a = a.Clone()
	a.Text = text
	t.s.doc.ReplaceAnnotation(a)
}

// sync shows the selected annotation's text.
func (t *textSection) sync() {
	if a, ok := t.s.selectedAnnot(); ok && a.Kind == annot.Text {
		t.field.SetText(a.Text)
	}
}

// edit focuses the field with all its text selected. Gio only gives the
// focus to a field it has seen laid out, so the request waits a frame.
func (t *textSection) edit() {
	t.sync()
	t.focusIn = 2
}

// View returns the section for this frame.
func (t *textSection) View() giode.Element {
	if t.focusIn > 0 {
		t.focusIn--
		if t.focusIn == 0 {
			t.field.Focus()
			t.field.SelectAll()
		} else {
			t.s.win.Post(func() {}) // one more frame
		}
	}
	return components.Section("Text", t.field)
}
