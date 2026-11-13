package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/internal/pagecache"
	"gopdf/ui/theme"
)

// thumbnailCard is one page in the thumbnail list: click to go there,
// right-click for page actions.
type thumbnailCard struct {
	s     *Studio
	page  int
	press *giode.PressableElement
}

func newThumbnailCard(s *Studio, page int) *thumbnailCard {
	c := &thumbnailCard{s: s, page: page}
	c.press = giode.Pressable(c.content).Label(fmt.Sprintf("Page %d", page+1)).
		OnClick(func(giode.ClickEvent) { s.goToPage(page) }).
		OnSecondaryClick(func(e giode.ClickEvent) { s.ui.pageMenu.open(page, e) })
	return c
}

// View returns the card for this frame.
func (c *thumbnailCard) View() giode.Element { return c.press }

func (c *thumbnailCard) content(st giode.PressState) giode.Element {
	info := c.s.doc.Doc().Pages[c.page]
	pw, ph := info.Size()
	w := theme.ThumbWidth
	h := int(float64(w) * ph / pw)
	active := c.page == c.s.view.current
	border, bg := theme.Border, "transparent"
	switch {
	case active:
		border, bg = theme.Primary, theme.Selected
	case st.Hovered:
		bg = theme.Hover
	}
	label := giode.Styles{FontSize: 11.5, Color: theme.Muted, FontWeight: giode.FontWeightMedium}
	if active {
		label.Color = theme.Primary
	}
	thumb := giode.Box(giode.Styles{
		Width: properties.Width(w), Height: properties.Height(h), Background: "#ffffff",
		BorderWidth: 1, BorderColor: border, BorderRadius: 4, Overflow: giode.OverflowHidden,
		BoxShadow: giode.BoxShadow{Y: 1, Blur: 3, Color: theme.PageShadow},
	}, c.image(w, h, pw))
	return giode.Stack(giode.Styles{
		Align: giode.AlignCenter, Gap: 6, Padding: giode.Inset{Top: 8, Bottom: 6, Left: 14, Right: 14},
		BorderRadius: 10, Background: bg,
	}, thumb, giode.Text(fmt.Sprint(c.page+1), label))
}

func (c *thumbnailCard) image(w, h int, pageWidth float64) giode.Element {
	dpi := int(float64(w) * float64(giode.DisplayScale()) * 72 / pageWidth)
	key := pagecache.Key{Version: c.s.doc.Version(), Page: c.page, DPI: max(dpi, 8)}
	img, fallback, ok := c.s.thumbs.Get(key)
	if !ok {
		img = fallback
	}
	if img == nil {
		return giode.Box(giode.Styles{Width: -1, Height: -1})
	}
	return giode.Image(img, giode.Styles{Width: properties.Width(w), Height: properties.Height(h), Fit: giode.FitStretch})
}
