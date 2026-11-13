package studio

// Zoom limits and steps.
const (
	minZoom = 0.25
	maxZoom = 4.0
)

var zoomSteps = []float64{0.25, 0.33, 0.5, 0.67, 0.75, 0.9, 1, 1.1, 1.25, 1.5, 1.75, 2, 2.5, 3, 4}

// sidebarTab is the left panel's content.
type sidebarTab int

const (
	tabPages sidebarTab = iota
	tabOutline
	tabSearch
)

// viewState is how the document is shown.
type viewState struct {
	zoom          float64
	fitWidth      bool
	current       int // page shown at the top of the viewer
	viewportWidth int // viewer width in dp, measured each frame
	sidebar       bool
	inspector     bool
	tab           sidebarTab
}

func newViewState() viewState {
	return viewState{zoom: 1, fitWidth: true, sidebar: true, inspector: true}
}

// nextZoom returns the zoom step after (dir > 0) or before z.
func nextZoom(z float64, dir int) float64 {
	if dir > 0 {
		for _, s := range zoomSteps {
			if s > z+0.001 {
				return s
			}
		}
		return maxZoom
	}
	for i := len(zoomSteps) - 1; i >= 0; i-- {
		if zoomSteps[i] < z-0.001 {
			return zoomSteps[i]
		}
	}
	return minZoom
}
