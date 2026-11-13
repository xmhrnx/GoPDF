package studio

// fitMargin is the space kept on each side of a page fitted to the width.
const fitMargin = 40

// zoomBy steps the zoom in (dir > 0) or out.
func (s *Studio) zoomBy(dir int) {
	s.setZoom(nextZoom(s.effectiveZoom(), dir))
}

func (s *Studio) setZoom(z float64) {
	s.view.fitWidth = false
	s.view.zoom = min(max(z, minZoom), maxZoom)
}

func (s *Studio) fitWidth() { s.view.fitWidth = true }

// effectiveZoom is the zoom in use: fixed, or fitting the widest page.
func (s *Studio) effectiveZoom() float64 {
	if !s.view.fitWidth || s.doc == nil || s.view.viewportWidth <= 0 {
		return s.view.zoom
	}
	widest := 0.0
	for _, p := range s.doc.Doc().Pages {
		w, _ := p.Size()
		widest = max(widest, w)
	}
	if widest == 0 {
		return 1
	}
	return min(max(float64(s.view.viewportWidth-2*fitMargin)/widest, minZoom), maxZoom)
}

// goToPage scrolls the viewer to page.
func (s *Studio) goToPage(page int) {
	if s.doc == nil || page < 0 || page >= s.pageCount() {
		return
	}
	s.view.current = page
	s.ui.viewer.scrollTo(page)
}

func (s *Studio) toggleSidebar()   { s.view.sidebar = !s.view.sidebar }
func (s *Studio) toggleInspector() { s.view.inspector = !s.view.inspector }

func (s *Studio) showTab(t sidebarTab) {
	s.view.sidebar = true
	s.view.tab = t
	s.ui.sidebar.syncTab()
}
