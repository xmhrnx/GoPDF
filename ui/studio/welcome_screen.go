package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// welcomeWidth is the width of the welcome column, in dp.
const welcomeWidth = 640

// welcomeScreen shows when no document is open: open, recent files and
// what the app can do.
type welcomeScreen struct {
	s      *Studio
	open   *giode.ButtonElement
	scroll *scroll.Scroll
}

func newWelcomeScreen(s *Studio) *welcomeScreen {
	w := &welcomeScreen{s: s}
	w.open = giode.Button("Open PDF…", giode.Styles{
		Background: theme.Primary, Color: theme.OnPrimary, BorderRadius: 10, FontSize: 14,
		Padding: giode.SymmetricInset(20, 11),
	}).Icon(theme.Icon("folder-open", 17, "")).OnClick(s.openDialog)
	w.scroll = giode.Scroll(w.content, giode.Styles{ScrollBar: theme.ScrollBar})
	return w
}

// View returns the screen for this frame.
func (w *welcomeScreen) View() giode.Element {
	return giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Background: theme.Background}, w.scroll)
}

func (w *welcomeScreen) content() giode.Element {
	return giode.Row(giode.Styles{Width: -1, Justify: giode.JustifyCenter, Padding: giode.SymmetricInset(24, 56)},
		giode.Stack(giode.Styles{Width: -1, MaxWidth: welcomeWidth, Gap: 36},
			w.hero(),
			w.features(),
			w.recent(),
		))
}

func (w *welcomeScreen) hero() giode.Element {
	return giode.Stack(giode.Styles{Width: -1, Align: giode.AlignCenter, Gap: 14},
		giode.Box(giode.Styles{BorderRadius: 18, BoxShadow: giode.BoxShadow{Y: 10, Blur: 30, Color: "#4f46e555"}}, logoBadge(64)),
		giode.Text("GoPDF", theme.Display),
		giode.Text("Annotate, organize, merge and protect your PDFs.", giode.Styles{FontSize: 15, Color: theme.Muted, TextAlign: giode.TextAlignCenter}),
		giode.Box(giode.Styles{Padding: giode.Inset{Top: 10}}, w.open),
		giode.Row(giode.Styles{Gap: 6, Align: giode.AlignCenter},
			giode.Text("or drop a file anywhere  ·", theme.Caption), components.Kbd("Mod+O")),
	)
}

func (w *welcomeScreen) features() giode.Element {
	return giode.Row(giode.Styles{Width: -1, Gap: 12, Wrap: giode.WrapWrap, Justify: giode.JustifyCenter, Align: giode.AlignStretch},
		featureCard("pen-line", "Annotate", "Text, ink, shapes, highlights"),
		featureCard("layout-grid", "Organize", "Rotate, reorder, duplicate"),
		featureCard("files", "Merge & split", "Combine or extract pages"),
		featureCard("shield-check", "Protect", "AES-256 passwords"),
	)
}

func (w *welcomeScreen) recent() giode.Element {
	list := w.s.recents.list
	if len(list) == 0 {
		return giode.Box(giode.Styles{})
	}
	rows := []giode.Element{giode.Text("Recent", theme.Overline)}
	for _, e := range list {
		rows = append(rows, giode.Keyed("recent", e.Path, func() *recentRow { return newRecentRow(w.s, e) }).View())
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 4}, rows...)
}
