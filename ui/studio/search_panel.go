package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/input"

	"gopdf/internal/pdf"
	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// maxMatches caps the results of one search.
const maxMatches = 500

// searchPanel finds text across the document.
type searchPanel struct {
	s       *Studio
	field   *input.Input
	list    *giode.LazyListElement
	query   string
	results []pdf.Match
	running bool
	gen     int // discards results of a search superseded by another
}

func newSearchPanel(s *Studio) *searchPanel {
	sp := &searchPanel{s: s}
	sp.field = components.NewTextInput("Search text…").Submit(sp.run)
	sp.list = giode.LazyList(func() int { return len(sp.results) }, func(i int) any { return i }, sp.row,
		giode.Styles{Gap: 2, ScrollBar: theme.ScrollBar})
	return sp
}

// reset clears results that refer to an older version of the document.
func (sp *searchPanel) reset() {
	sp.results, sp.query = nil, ""
	sp.gen++
	sp.running = false
}

// focus shows the panel and puts the caret in the field.
func (sp *searchPanel) focus() {
	sp.s.showTab(tabSearch)
	sp.field.Focus()
}

func (sp *searchPanel) run(query string) {
	if sp.s.renderer == nil || query == "" {
		return
	}
	sp.gen++
	gen, r, n := sp.gen, sp.s.renderer, sp.s.pageCount()
	sp.query, sp.results, sp.running = query, nil, true
	go func() {
		var found []pdf.Match
		for p := 0; p < n && len(found) < maxMatches; p++ {
			text, err := r.Text(p)
			if err != nil {
				break // the renderer was replaced; a newer search will run
			}
			found = append(found, pdf.FindInText(p, text, query, maxMatches-len(found))...)
		}
		sp.s.win.Post(func() {
			if gen == sp.gen {
				sp.results, sp.running = found, false
			}
		})
	}()
}

// View returns the panel for this frame.
func (sp *searchPanel) View() giode.Element {
	var body giode.Element
	switch {
	case sp.running:
		body = giode.Row(giode.Styles{Width: -1, Justify: giode.JustifyCenter, Padding: giode.UniformInset(24)},
			giode.Spinner(giode.Styles{Width: 20, Color: theme.Primary}))
	case sp.query == "":
		body = components.EmptyHint("search", "Search the document", "Type a word or phrase and press Enter.")
	case len(sp.results) == 0:
		body = components.EmptyHint("search", "No matches", fmt.Sprintf("Nothing found for “%s”.", sp.query))
	default:
		body = giode.Stack(giode.Styles{Width: -1, Height: -1},
			giode.Box(giode.Styles{Padding: giode.SymmetricInset(16, 4)},
				giode.Text(fmt.Sprintf("%d matches", len(sp.results)), theme.Caption)),
			giode.Box(giode.Styles{Width: -1, FlexGrow: 1}, sp.list),
		)
	}
	return giode.Stack(giode.Styles{Width: -1, Height: -1},
		giode.Box(giode.Styles{Width: -1, Padding: giode.SymmetricInset(10, 0)}, sp.field),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Padding: giode.Inset{Top: 8}}, body),
	)
}

func (sp *searchPanel) row(i int) giode.Element {
	return giode.Box(giode.Styles{Width: -1, Padding: giode.SymmetricInset(8, 0)}, sp.result(i))
}

func (sp *searchPanel) result(i int) giode.Element {
	return giode.Keyed("match", i, func() *giode.PressableElement {
		return giode.Pressable(func(st giode.PressState) giode.Element {
			if i >= len(sp.results) {
				return giode.Box(giode.Styles{})
			}
			m := sp.results[i]
			bg := "transparent"
			if st.Hovered {
				bg = theme.Hover
			}
			return giode.Stack(giode.Styles{Width: -1, Gap: 3, Padding: giode.SymmetricInset(10, 7), BorderRadius: 8, Background: bg},
				giode.Text(fmt.Sprintf("Page %d", m.Page+1), giode.Styles{FontSize: 11, Color: theme.Primary, FontWeight: giode.FontWeightSemiBold}),
				giode.Text(m.Snippet, giode.Styles{FontSize: 12, Color: theme.TextColor, MaxLines: 3}),
			)
		}).OnClick(func(giode.ClickEvent) {
			if i < len(sp.results) {
				sp.s.goToPage(sp.results[i].Page)
			}
		})
	})
}
