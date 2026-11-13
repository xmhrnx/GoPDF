package studio

import (
	"strings"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/internal/signature"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// maxNameLength caps a typed signature, in characters.
const maxNameLength = 60

// typePanel writes a typed name in a script typeface.
type typePanel struct {
	s     *Studio
	name  *input.Input
	font  int
	cards []*fontCard
}

func newTypePanel(s *Studio) *typePanel {
	t := &typePanel{s: s}
	t.name = components.NewTextInput("Your full name").MaxLength(maxNameLength)
	for i := range s.scripts {
		t.cards = append(t.cards, newFontCard(t, i))
	}
	return t
}

// reset empties the name and selects the typeface used last.
func (t *typePanel) reset() {
	t.name.SetText("")
	t.font = 0
	for i, f := range t.s.scripts {
		if f.Name == t.s.memory.prefs.SignatureFont {
			t.font = i
		}
	}
}

func (t *typePanel) pick(i int) { t.font = i }

// previewText is what the font cards show: the name, or a sample.
func (t *typePanel) previewText() string {
	if name := strings.TrimSpace(t.name.Text()); name != "" {
		return name
	}
	return "Your Name"
}

// signature writes the name in the chosen typeface.
func (t *typePanel) signature() (signature.Signature, error) {
	return signature.FromText(t.name.Text(), t.s.scripts[t.font].Data)
}

// fontName is the chosen typeface's name.
func (t *typePanel) fontName() string { return t.s.scripts[t.font].Name }

// View returns the panel for this frame.
func (t *typePanel) View() giode.Element {
	cards := make([]giode.Element, len(t.cards))
	for i, c := range t.cards {
		cards[i] = c.View()
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 14},
		components.Field("Name", t.name),
		components.Field("Style", giode.Row(giode.Styles{Width: -1, Gap: 8, Wrap: giode.WrapWrap}, cards...)),
		giode.Text("Saved as outlines, so it looks the same in every PDF viewer.", theme.Caption),
	)
}
