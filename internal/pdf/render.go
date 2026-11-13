package pdf

import (
	"errors"
	"image"
	"strings"
	"sync"

	"github.com/gen2brain/go-fitz"
)

// mupdfStore caps MuPDF's resource cache (fonts, decoded images) per
// document, in bytes; its 256 MB default makes a viewer heavy.
const mupdfStore = 48 << 20

func init() { fitz.MaxStore = mupdfStore }

// ErrClosed is returned by a Renderer used after Close.
var ErrClosed = errors.New("renderer closed")

// Renderer rasterizes pages and extracts text with MuPDF. It is safe for
// use from several goroutines; calls are serialized, and calls after
// Close fail with ErrClosed instead of touching freed memory.
type Renderer struct {
	mu  sync.Mutex
	doc *fitz.Document
}

// OutlineItem is one entry of the document's table of contents.
type OutlineItem struct {
	Title string
	Page  int // zero-based
	Level int // 1 for top level
}

// NewRenderer opens data (unencrypted PDF bytes, which must not change
// afterwards) for rendering.
func NewRenderer(data []byte) (*Renderer, error) {
	doc, err := fitz.NewFromMemory(data)
	if err != nil {
		return nil, err
	}
	return &Renderer{doc: doc}, nil
}

// Render draws page (zero-based) at dpi dots per inch.
func (r *Renderer) Render(page int, dpi float64) (*image.RGBA, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.doc == nil {
		return nil, ErrClosed
	}
	return r.doc.ImageDPI(page, dpi)
}

// Text returns the plain text of page.
func (r *Renderer) Text(page int) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.doc == nil {
		return "", ErrClosed
	}
	return r.doc.Text(page)
}

// Outline returns the table of contents, or nil when there is none.
func (r *Renderer) Outline() []OutlineItem {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.doc == nil {
		return nil
	}
	toc, err := r.doc.ToC()
	if err != nil {
		return nil
	}
	items := make([]OutlineItem, 0, len(toc))
	for _, o := range toc {
		if t := strings.TrimSpace(o.Title); t != "" {
			items = append(items, OutlineItem{Title: t, Page: o.Page, Level: o.Level})
		}
	}
	return items
}

// Close releases the MuPDF document.
func (r *Renderer) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.doc == nil {
		return nil
	}
	err := r.doc.Close()
	r.doc = nil
	return err
}
