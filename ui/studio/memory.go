package studio

import (
	"maps"
	"sync"
	"time"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/settings"

	"gopdf/internal/annot"
	"gopdf/internal/prefs"
)

// saveDelay gathers quick changes (a color dragged in the picker) into
// one write.
const saveDelay = 400 * time.Millisecond

// memory is what the app remembers between runs: the last style used
// for each kind of markup, and the signature and ID to use.
type memory struct {
	prefs prefs.Prefs

	mu      sync.Mutex
	pending *prefs.Prefs // the copy the next write saves
	timer   *time.Timer
}

func newMemory() *memory { return &memory{} }

// load reads the saved choices on a background goroutine.
func (m *memory) load(s *Studio) {
	go func() {
		p := settings.Value(keyPrefs, prefs.Prefs{})
		s.win.Post(func() { m.prefs = p })
	}()
}

// style returns the remembered style of kind k.
func (m *memory) style(k annot.Kind) annot.Style { return m.prefs.StyleFor(k) }

// setStyle remembers st for k and saves it.
func (m *memory) setStyle(s *Studio, k annot.Kind, st annot.Style) {
	m.prefs.SetStyle(k, st)
	m.save(s)
}

// setSignature remembers the signature to place.
func (m *memory) setSignature(s *Studio, id string) {
	m.prefs.Signature = id
	m.save(s)
}

// save writes the prefs shortly, once changes stop coming.
func (m *memory) save(s *Studio) {
	p := m.prefs
	p.Styles = maps.Clone(m.prefs.Styles) // the UI goroutine keeps changing the original
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pending = &p
	if m.timer != nil {
		m.timer.Reset(saveDelay)
		return
	}
	m.timer = time.AfterFunc(saveDelay, func() { m.flush(s) })
}

func (m *memory) flush(s *Studio) {
	m.mu.Lock()
	p := m.pending
	m.pending = nil
	m.mu.Unlock()
	if p == nil {
		return
	}
	if err := settings.Save(keyPrefs, *p); err != nil {
		s.win.Post(func() { s.toasts.Show(giode.ToastWarning, "Couldn't save your settings: "+err.Error()) })
	}
}
