// Package studio is the GoPDF interface: the window's views, the
// state they show and the actions they run.
package studio

import (
	"image"
	"net/url"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/toaster"
	"github.com/xmhrnx/giode/styles/properties"

	"gopdf/internal/annot"
	"gopdf/internal/editor"
	"gopdf/internal/pagecache"
	"gopdf/internal/pdf"
	"gopdf/ui/theme"
)

// Cache sizes: rendered pages are large, thumbnails small.
const (
	pageCacheLimit  = 10
	thumbCacheLimit = 240
)

// Studio is the application: it owns the open document and every view.
type Studio struct {
	win     *giode.App
	logo    image.Image
	scripts []theme.ScriptFont // typefaces for typed signatures
	recents *recentFiles
	toasts  *toaster.Toaster

	docs     []*openDoc      // the open documents, in tab order
	active   int             // index of the active one in docs, -1 for none
	doc      *editor.Session // the active document; nil while none is open
	renderer *pdf.Renderer
	rendered int // session version the renderer shows
	pages    *pagecache.Cache
	thumbs   *pagecache.Cache
	outline  []pdf.OutlineItem
	busy     string // what runs in the background, "" when idle

	tool        Tool
	memory      *memory        // per-kind styles and the active signature, kept on disk
	signatures  *signatureBook // saved signatures
	ids         *digitalIDs    // digital IDs to sign with
	clipboard   *clipboard     // copied markup
	nudging     nudgeState     // groups arrow-key moves into one undo step
	live        liveGroup      // groups continuous style edits into one undo step
	selected    int64          // selected annotation ID, 0 for none
	skipSealing bool           // place the next signature without a digital ID
	drawing     *annot.Annotation
	moving      *moveGesture
	view        viewState
	dragOver    bool

	ui *views
}

// New returns the studio, showing logo as the app's mark and offering
// scripts for typed signatures; Attach it to a window before running it.
func New(logo image.Image, scripts []theme.ScriptFont) *Studio {
	s := &Studio{logo: logo, scripts: scripts, active: -1, tool: ToolSelect, view: newViewState(), clipboard: &clipboard{}}
	s.memory = newMemory()
	s.signatures = newSignatureBook()
	s.ids = newDigitalIDs()
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
	s.memory.load(s)
	s.signatures.load(s)
	s.ids.load(s)
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
			s.openPath(p)
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
// window open while documents have unsaved changes, and asks about each.
func (s *Studio) OnClose() bool {
	if !s.anyDirty() {
		return true
	}
	s.closeAll(s.win.Close)
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
		s.openPath(u.Path)
	}
}
