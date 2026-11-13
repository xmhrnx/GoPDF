// Package studio is the GoPDF interface: the window's views, the
// state they show and the actions they run.
package studio

import (
	"net/url"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/toaster"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/internal/annot"
	"gopdf/internal/editor"
	"gopdf/internal/pagecache"
	"gopdf/internal/pdf"
)

// Cache sizes: rendered pages are large, thumbnails small.
const (
	pageCacheLimit  = 10
	thumbCacheLimit = 240
)

// Studio is the application: it owns the open document and every view.
type Studio struct {
	win     *giode.App
	recents *recentFiles
	toasts  *toaster.Toaster

	doc      *editor.Session // nil while no document is open
	renderer *pdf.Renderer
	rendered int // session version the renderer shows
	pages    *pagecache.Cache
	thumbs   *pagecache.Cache
	outline  []pdf.OutlineItem
	busy     string // what runs in the background, "" when idle

	tool     Tool
	styles   map[Tool]annot.Style // last style used with each tool
	selected int64                // selected annotation ID, 0 for none
	drawing  *annot.Annotation
	moving   *moveGesture
	view     viewState
	dragOver bool

	ui *views
}

// New returns the studio; Attach it to a window before running it.
func New() *Studio {
	s := &Studio{tool: ToolSelect, styles: defaultToolStyles(), view: newViewState()}
	s.toasts = giode.Toaster().Placement(giode.ToastBottomEnd).Max(3).
		Inset(giode.Inset{Bottom: 36, Right: 8})
	s.recents = newRecentFiles()
	return s
}

// Attach connects the studio to its window and builds the views.
func (s *Studio) Attach(win *giode.App) {
	s.win = win
	invalidate := func() { win.Post(func() {}) }
	s.pages = pagecache.New(pageCacheLimit, invalidate, properties.ForgetImage)
	s.thumbs = pagecache.New(thumbCacheLimit, invalidate, properties.ForgetImage)
	s.ui = newViews(s)
	s.recents.load(s)
}

// View is the window's view function.
func (s *Studio) View() giode.Element {
	s.syncRenderer()
	return s.ui.root()
}

// OnFileDrop opens the first PDF dropped onto the window.
func (s *Studio) OnFileDrop(paths []string) {
	s.dragOver = false
	for _, p := range paths {
		if isPDF(p) {
			s.guardUnsaved(func() { s.openPath(p) })
			return
		}
	}
	s.toasts.Show(giode.ToastWarning, "Drop a PDF file to open it")
}

// OnFileDrag shows the drop overlay while files hover over the window.
func (s *Studio) OnFileDrag(d giode.FileDrag) {
	s.dragOver = d.Kind == giode.FileDragEnter || d.Kind == giode.FileDragOver
}

// OnClose runs when the title bar's close button is pressed; it keeps the
// window open while there are unsaved changes and asks first.
func (s *Studio) OnClose() bool {
	if s.doc == nil || !s.doc.Dirty() {
		return true
	}
	s.guardUnsaved(s.win.Close)
	return false
}

// hasDoc reports whether a document is open.
func (s *Studio) hasDoc() bool { return s.doc != nil }

// idle reports whether no background work is running.
func (s *Studio) idle() bool { return s.busy == "" }

// canEdit reports whether a document is open and nothing is running.
func (s *Studio) canEdit() bool { return s.hasDoc() && s.idle() }

// OpenFiles opens the first PDF among paths (command-line arguments).
func (s *Studio) OpenFiles(paths []string) {
	for _, p := range paths {
		if isPDF(p) {
			s.win.Post(func() { s.openPath(p) })
			return
		}
	}
}

// OnOpenURL opens a PDF handed over as a file:// URL (macOS "Open with").
func (s *Studio) OnOpenURL(u *url.URL) {
	if u.Scheme == "file" && isPDF(u.Path) {
		s.guardUnsaved(func() { s.openPath(u.Path) })
	}
}
