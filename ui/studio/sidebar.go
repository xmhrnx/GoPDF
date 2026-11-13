package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/segmented"

	"gopdf/ui/components"
	"gopdf/ui/theme"
)

// sidebar is the left panel: page thumbnails, outline or search.
type sidebar struct {
	s       *Studio
	tabs    *segmented.Control
	thumbs  *thumbnailList
	outline *outlinePanel
}

func newSidebar(s *Studio) *sidebar {
	sb := &sidebar{s: s, thumbs: newThumbnailList(s), outline: newOutlinePanel(s)}
	sb.tabs = components.NewSegmented("Pages", "Outline", "Find").
		OnChange(func(i int) { s.view.tab = sidebarTab(i) })
	return sb
}

// syncTab shows the tab chosen elsewhere (a shortcut).
func (sb *sidebar) syncTab() { sb.tabs.SetSelected(int(sb.s.view.tab)) }

// View returns the panel for this frame.
func (sb *sidebar) View() giode.Element {
	var body giode.Element
	switch sb.s.view.tab {
	case tabOutline:
		body = sb.outline.View()
	case tabSearch:
		body = sb.s.ui.search.View()
	default:
		body = sb.thumbs.View()
	}
	panel := giode.Stack(giode.Styles{Width: theme.SidebarWidth, Height: -1, Background: theme.Surface},
		giode.Box(giode.Styles{Width: -1, Padding: giode.UniformInset(10)}, sb.tabs),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1}, body),
	)
	return giode.Row(giode.Styles{Height: -1, Wrap: giode.WrapNoWrap}, panel, components.VRule())
}
