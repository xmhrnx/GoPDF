package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/internal/pdf"
)

// signatureFieldLayer makes the document's signed fields on page
// clickable with the select tool, showing who signed.
func signatureFieldLayer(s *Studio, page int, z float64) []giode.Element {
	if s.tool != ToolSelect {
		return nil
	}
	var out []giode.Element
	all := s.doc.Doc().Signatures
	for i, f := range all {
		if f.Page != page || f.Rect.Dx() == 0 {
			continue
		}
		key := [3]int{s.doc.Version(), page, i}
		field := giode.Keyed("signature-field", key, func() *giode.PressableElement {
			return newSignatureField(s, f, all)
		})
		field.Styles(giode.Styles{
			Width: properties.Width(int(f.Rect.Dx() * z)), Height: properties.Height(int(f.Rect.Dy() * z)),
		})
		out = append(out, giode.Box(giode.Styles{
			Position: giode.PositionAbsolute,
			Left:     properties.Left(edge(f.Rect.Min.X * z)), Top: properties.Top(edge(f.Rect.Min.Y * z)),
		}, field))
	}
	return out
}

// newSignatureField is the clickable area over one signed field.
func newSignatureField(s *Studio, f pdf.SignatureField, all []pdf.SignatureField) *giode.PressableElement {
	level := statusOf(f, all).level
	return giode.Pressable(func(st giode.PressState) giode.Element {
		if !st.Hovered {
			return giode.Box(giode.Styles{Width: -1, Height: -1})
		}
		return giode.Box(giode.Styles{
			Width: -1, Height: -1, BorderWidth: 2, BorderRadius: 4, BorderColor: sealColor(level),
		})
	}).Label("Signature by " + f.Name).OnClick(func(giode.ClickEvent) {
		s.ui.sigDetails.show([]pdf.SignatureField{f})
	})
}
