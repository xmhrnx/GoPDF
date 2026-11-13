package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/ui/theme"
)

// Font card size and preview text size.
const (
	fontCardWidth  = 202
	fontCardHeight = 70
	fontPreview    = 28
	previewHeight  = 42
)

// fontCard previews the typed name in one script typeface; clicking
// picks the typeface.
type fontCard struct {
	panel *typePanel
	index int
	press *giode.PressableElement
}

func newFontCard(panel *typePanel, index int) *fontCard {
	c := &fontCard{panel: panel, index: index}
	c.press = giode.Pressable(c.content).Label(panel.s.scripts[index].Name).
		OnClick(func(giode.ClickEvent) { panel.pick(index) })
	return c
}

// View returns the card for this frame.
func (c *fontCard) View() giode.Element { return c.press }

func (c *fontCard) content(st giode.PressState) giode.Element {
	font := c.panel.s.scripts[c.index]
	border, bg := theme.Border, "#ffffff"
	switch {
	case c.panel.font == c.index:
		border, bg = theme.Primary, giode.Mix("theme.primary", "#ffffff", 0.06)
	case st.Hovered:
		border = theme.Faint
	}
	return giode.Stack(giode.Styles{
		Width: fontCardWidth, Height: fontCardHeight, Background: bg, BorderRadius: 10,
		BorderWidth: 1, BorderColor: border, Padding: giode.Inset{Top: 4, Left: 12, Right: 12, Bottom: 6},
		Justify: giode.JustifySpaceBetween, Overflow: giode.OverflowHidden,
	},
		giode.Box(giode.Styles{Width: -1, Height: previewHeight, Overflow: giode.OverflowHidden, Align: giode.AlignCenter},
			giode.Text(c.panel.previewText(), giode.Styles{
				FontFamily: properties.FontFamily(font.Name), FontSize: fontPreview, Color: "#1e293b",
				MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis,
			})),
		giode.Text(font.Name, giode.Styles{FontSize: 10.5, Color: "#64748b"}),
	)
}
