package studio

import (
	"fmt"
	"math"

	"github.com/xmhrnx/giode"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// zoomControl shows the zoom level with buttons around it.
type zoomControl struct {
	s     *Studio
	out   *components.ToolButton
	in    *components.ToolButton
	fit   *components.ToolButton
	level *giode.PressableElement
}

func newZoomControl(s *Studio) *zoomControl {
	z := &zoomControl{s: s}
	z.out = components.NewToolButton("zoom-out", "Zoom out", func() { s.zoomBy(-1) })
	z.in = components.NewToolButton("zoom-in", "Zoom in", func() { s.zoomBy(1) })
	z.fit = components.NewToolButton("move-horizontal", "Fit width", s.fitWidth).
		Active(func() bool { return s.view.fitWidth })
	z.level = giode.Pressable(z.levelContent).Label("Actual size").
		OnClick(func(giode.ClickEvent) { s.setZoom(1) })
	return z
}

// View returns the control for this frame.
func (z *zoomControl) View() giode.Element {
	return giode.Row(giode.Styles{Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		z.out.View(), z.level, z.in.View(), components.VDivider(20), z.fit.View(),
	)
}

func (z *zoomControl) levelContent(st giode.PressState) giode.Element {
	bg := "transparent"
	if st.Hovered {
		bg = theme.Hover
	}
	pct := fmt.Sprintf("%d%%", int(math.Round(z.s.effectiveZoom()*100)))
	return giode.Box(giode.Styles{
		Width: 56, Height: 28, BorderRadius: 6, Background: bg, Align: giode.AlignCenter, Justify: giode.JustifyCenter,
	}, giode.Text(pct, giode.Styles{FontSize: 12.5, Color: theme.TextColor, FontWeight: giode.FontWeightMedium}))
}
