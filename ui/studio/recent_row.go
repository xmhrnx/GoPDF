package studio

import (
	"fmt"
	"path/filepath"

	"github.com/xmhrnx/giode"

	"gopdf/internal/recent"
	"gopdf/ui/theme"
)

// recentRow is one recently opened file on the welcome screen.
type recentRow struct {
	entry recent.Entry
	press *giode.PressableElement
}

func newRecentRow(s *Studio, e recent.Entry) *recentRow {
	r := &recentRow{entry: e}
	r.press = giode.Pressable(r.content).Label(e.Name()).
		OnClick(func(giode.ClickEvent) { s.openPath(e.Path) })
	return r
}

// View returns the row for this frame.
func (r *recentRow) View() giode.Element { return r.press }

func (r *recentRow) content(st giode.PressState) giode.Element {
	bg, border := theme.Surface, theme.Border
	if st.Hovered {
		bg, border = theme.Hover, theme.Primary
	}
	meta := filepath.Dir(r.entry.Path)
	if r.entry.Pages > 0 {
		meta = fmt.Sprintf("%d pages · %s", r.entry.Pages, meta)
	}
	return giode.Row(giode.Styles{
		Width: -1, Gap: 12, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(14, 10),
		BorderRadius: 10, Background: bg, BorderWidth: 1, BorderColor: border,
	},
		giode.Box(giode.Styles{Width: 32, Height: 32, BorderRadius: 8, Background: theme.Selected, Align: giode.AlignCenter, Justify: giode.JustifyCenter},
			theme.Icon("file-text", 16, theme.Primary)),
		giode.Stack(giode.Styles{Gap: 2, FlexGrow: 1},
			giode.Text(r.entry.Name(), giode.Styles{FontSize: 13, Color: theme.TextColor, FontWeight: giode.FontWeightMedium, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Text(meta, giode.Styles{FontSize: 11.5, Color: theme.Muted, TextOverflow: giode.TextOverflowEllipsis}),
		),
		giode.Text(r.entry.Opened.Format("Jan 2"), theme.Caption),
	)
}
