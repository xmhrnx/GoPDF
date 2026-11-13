package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// statusBar is the strip at the bottom: page navigation and state.
type statusBar struct {
	s    *Studio
	prev *components.ToolButton
	next *components.ToolButton
}

func newStatusBar(s *Studio) *statusBar {
	return &statusBar{
		s:    s,
		prev: components.NewToolButton("chevron-left", "Previous page", func() { s.goToPage(s.view.current - 1) }).Compact(),
		next: components.NewToolButton("chevron-right", "Next page", func() { s.goToPage(s.view.current + 1) }).Compact(),
	}
}

// View returns the bar for this frame.
func (b *statusBar) View() giode.Element {
	n := b.s.pageCount()
	nav := giode.Row(giode.Styles{Gap: 4, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		b.prev.View(),
		giode.Text(fmt.Sprintf("Page %d of %d", b.s.view.current+1, n), theme.Caption),
		b.next.View(),
	)
	right := []giode.Element{}
	if b.s.busy != "" {
		right = append(right, giode.Spinner(giode.Styles{Width: 12, Color: theme.Primary}), giode.Text(b.s.busy+"…", theme.Caption))
	} else if b.s.doc.Dirty() {
		right = append(right, dot(theme.Primary), giode.Text("Unsaved changes", theme.Caption))
	} else {
		right = append(right, dot(theme.Success), giode.Text("Saved", theme.Caption))
	}
	if b.s.doc.Security.Enabled() {
		right = append(right, theme.Icon("lock", 12, theme.Muted))
	}
	bar := giode.Row(giode.Styles{
		Width: -1, Height: theme.StatusBarHeight, Align: giode.AlignCenter, Justify: giode.JustifySpaceBetween,
		Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(8, 0), Background: theme.Surface,
	}, nav, giode.Row(giode.Styles{Gap: 6, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.Inset{Right: 8}}, right...))
	return giode.Stack(giode.Styles{Width: -1}, components.HRule(), bar)
}

// dot is a small status light.
func dot(color string) giode.Element {
	return giode.Box(giode.Styles{Width: 7, Height: 7, BorderRadius: 99, Background: color})
}
