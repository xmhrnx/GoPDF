package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	"gopdf/internal/annot"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// inspector is the right panel: the style of the active tool or the
// selected markup, its text, and the document's details.
type inspector struct {
	s      *Studio
	style  *styleSection
	sigs   *signatureSection
	text   *textSection
	doc    *documentSection
	scroll *scroll.Scroll
}

func newInspector(s *Studio) *inspector {
	in := &inspector{s: s, style: newStyleSection(s), sigs: newSignatureSection(s), text: newTextSection(s), doc: newDocumentSection(s)}
	in.scroll = giode.Scroll(in.content, giode.Styles{ScrollBar: theme.ScrollBar})
	return in
}

// syncText loads the selected text annotation into the text field.
func (in *inspector) syncText() { in.text.sync() }

// editText focuses the text field with its content selected.
func (in *inspector) editText() { in.text.edit() }

// View returns the panel for this frame.
func (in *inspector) View() giode.Element {
	panel := giode.Box(giode.Styles{Width: theme.InspectorWidth, Height: -1, Background: theme.Surface}, in.scroll)
	return giode.Row(giode.Styles{Height: -1, Wrap: giode.WrapNoWrap}, components.VRule(), panel)
}

func (in *inspector) content() giode.Element {
	kind, ok := in.subject()
	kids := []giode.Element{in.header(kind, ok)}
	if ok {
		kids = append(kids, in.style.View(kind))
		if a, sel := in.s.selectedAnnot(); sel && a.Kind == annot.Text {
			kids = append(kids, in.text.View())
		}
		if kind == annot.Signature {
			kids = append(kids, in.sigs.View())
		}
	}
	kids = append(kids, components.HRule(), in.doc.View())
	return giode.Stack(giode.Styles{Width: -1}, kids...)
}

// subject is the kind the style controls act on.
func (in *inspector) subject() (annot.Kind, bool) { return in.s.styledKind() }

func (in *inspector) header(kind annot.Kind, ok bool) giode.Element {
	title, hint := "Select", "Click markup to edit it, or pick a tool to draw."
	if _, sel := in.s.selectedAnnot(); sel {
		title, hint = kind.String(), "Selected · drag to move, Delete to remove"
	} else if ok {
		title, hint = kind.String(), "Drag on a page to draw"
		switch kind {
		case annot.Text:
			hint = "Click on a page to add text"
		case annot.Signature:
			hint = "Click on a page to place your signature"
		}
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 4, Padding: giode.Inset{Top: 18, Left: 16, Right: 16, Bottom: 6}},
		giode.Text(title, theme.Heading),
		giode.Text(hint, theme.Caption),
	)
}
