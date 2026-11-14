package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dropdown"
	"github.com/xmhrnx/giode/components/input"
	"github.com/xmhrnx/giode/components/toggle"

	"gopdf/internal/annot"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// sealSection makes the selected signature a digital signature: who
// signs, and why and where.
type sealSection struct {
	s        *Studio
	on       *toggle.Switch
	id       *dropdown.Select
	reason   *input.Input
	location *input.Input
	manage   *giode.ButtonElement
}

func newSealSection(s *Studio) *sealSection {
	sec := &sealSection{s: s}
	sec.on = giode.Input.Switch("Sign with a digital ID", giode.Styles{FontSize: 13, Color: theme.TextColor}).OnChange(s.setSealed)
	sec.id = giode.Input.Select(giode.Styles{Background: theme.Surface, Color: theme.TextColor, BorderColor: theme.Border, BorderRadius: theme.Radius, Width: -1}).
		OnChange(func(i int) {
			if i < len(s.ids.list) {
				s.sealWith(s.ids.list[i].Fingerprint)
			}
		})
	sec.reason = components.NewTextInput("e.g. I approve this document").OnChange(s.setSealReason)
	sec.location = components.NewTextInput("e.g. Berlin").OnChange(s.setSealLocation)
	sec.manage = components.NewSecondaryButton("Digital IDs…", func() { s.ui.digitalIDs.show(nil) }).Icon(theme.Icon("badge-check", 15, ""))
	return sec
}

// focused reports whether one of the section's fields has focus.
func (sec *sealSection) focused() bool { return sec.reason.Focused() || sec.location.Focused() }

// blur takes the focus from the section's fields.
func (sec *sealSection) blur() {
	sec.reason.Blur()
	sec.location.Blur()
}

// View returns the section for the selected signature a.
func (sec *sealSection) View(a annot.Annotation) giode.Element {
	sec.on.SetChecked(a.Seal != nil)
	kids := []giode.Element{sec.on}
	if a.Seal == nil {
		kids = append(kids, giode.Text("A digital signature seals the document with your digital ID: readers can click it to see who signed, and whether anything changed since.", theme.Caption))
		return components.Section("Digital signature", kids...)
	}
	sec.syncIDs(a.Seal.ID)
	syncText(sec.reason, a.Seal.Reason)
	syncText(sec.location, a.Seal.Location)
	kids = append(kids,
		components.Field("Sign as", sec.id),
		components.Field("Reason", sec.reason),
		components.Field("Location", sec.location),
		giode.Text("Applied when you save. Changing the document after that invalidates the signature.", theme.Caption),
		sec.manage,
	)
	return components.Section("Digital signature", kids...)
}

// syncIDs lists the stored IDs in the menu, with id selected.
func (sec *sealSection) syncIDs(id string) {
	labels := make([]string, len(sec.s.ids.list))
	selected := -1
	for i, info := range sec.s.ids.list {
		labels[i] = info.Identity.Name
		if info.Identity.Email != "" {
			labels[i] += " · " + info.Identity.Email
		}
		if info.Fingerprint == id {
			selected = i
		}
	}
	if selected < 0 {
		labels = append(labels, "Deleted digital ID")
		selected = len(labels) - 1
	}
	sec.id.SetOptions(labels...)
	sec.id.SetSelected(selected)
}

// syncText shows text in field unless the user is typing in it.
func syncText(field *input.Input, text string) {
	if !field.Focused() && field.Text() != text {
		field.SetText(text)
	}
}
